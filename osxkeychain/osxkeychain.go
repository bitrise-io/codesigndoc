package osxkeychain

import (
	"crypto/x509"
	"errors"
	"fmt"
	"unsafe"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-xcode/certificateutil"
)

/*
#cgo CFLAGS: -mmacosx-version-min=10.7 -D__MAC_OS_X_VERSION_MAX_ALLOWED=1060
#cgo LDFLAGS: -framework CoreFoundation -framework Security
#include <stdlib.h>
#include <CoreFoundation/CoreFoundation.h>
#include <Security/Security.h>
*/
import "C"

// ExportFromKeychain ...
func ExportFromKeychain(itemRefsToExport []C.CFTypeRef, isAskForPassword bool) ([]byte, error) {
	passphraseCString := C.CString("")
	defer C.free(unsafe.Pointer(passphraseCString))

	var exportedData C.CFDataRef
	var exportParams C.SecItemImportExportKeyParameters
	exportParams.keyUsage = 0
	exportParams.keyAttributes = 0
	exportParams.version = C.SEC_KEY_IMPORT_EXPORT_PARAMS_VERSION
	if isAskForPassword {
		exportParams.flags = C.kSecKeySecurePassphrase
		exportParams.passphrase = 0
		exportParams.alertTitle = 0

		promptText := C.CString("Enter a password which will be used to protect the exported items")
		defer C.free(unsafe.Pointer(promptText))
		exportParams.alertPrompt = convertCStringToCFString(promptText)
	} else {
		exportParams.flags = 0
		exportParams.passphrase = (C.CFTypeRef)(convertCStringToCFString(passphraseCString))
		exportParams.alertTitle = 0
		exportParams.alertPrompt = 0
	}

	// create a C array from the input
	ptr := (*unsafe.Pointer)(unsafe.Pointer(&itemRefsToExport[0]))
	cfArrayForExport := C.CFArrayCreate(
		C.kCFAllocatorDefault,
		ptr,
		C.CFIndex(len(itemRefsToExport)),
		&C.kCFTypeArrayCallBacks)

	// do the export!
	status := C.SecItemExport(C.CFTypeRef(cfArrayForExport),
		C.kSecFormatPKCS12,
		0, //C.kSecItemPemArmour, // Use kSecItemPemArmour to add PEM armour - the .p12 generated by Keychain Access.app does NOT have PEM armour
		&exportParams,
		&exportedData)

	if status != C.errSecSuccess {
		return nil, fmt.Errorf("SecItemExport: error (OSStatus): %d", status)
	}
	// exportedData now contains your PKCS12 data
	//  make sure it'll be released properly!
	defer C.CFRelease(C.CFTypeRef(exportedData))

	dataBytes := convertCFDataRefToGoBytes(exportedData)
	if dataBytes == nil || len(dataBytes) < 1 {
		return nil, errors.New("ExportFromKeychain: failed to convert export data - nil or empty")
	}
	log.Debugf("Export - success")

	return dataBytes, nil
}

func convertCFDataRefToGoBytes(cfdata C.CFDataRef) []byte {
	return C.GoBytes(unsafe.Pointer(C.CFDataGetBytePtr(cfdata)), (C.int)(C.CFDataGetLength(cfdata)))
}

// ReleaseRef ...
func ReleaseRef(refItem C.CFTypeRef) {
	C.CFRelease(refItem)
}

// ReleaseRefList ...
func ReleaseRefList(refItems []C.CFTypeRef) {
	for _, itm := range refItems {
		ReleaseRef(itm)
	}
}

// ReleaseIdentityWithRefList ...
func ReleaseIdentityWithRefList(refItems []IdentityWithRefModel) {
	for _, itm := range refItems {
		ReleaseRef(itm.KeychainRef)
	}
}

// CreateEmptyCFTypeRefSlice ...
func CreateEmptyCFTypeRefSlice() []C.CFTypeRef {
	return []C.CFTypeRef{}
}

// GetCertificateDataFromIdentityRef ...
func GetCertificateDataFromIdentityRef(identityRef C.CFTypeRef) (*x509.Certificate, error) {
	secIdentityRef := C.SecIdentityRef(identityRef)
	var secCertificateRef C.SecCertificateRef
	osStatusCode := C.SecIdentityCopyCertificate(secIdentityRef, &secCertificateRef)
	if osStatusCode != C.errSecSuccess {
		return nil, fmt.Errorf("Failed to call SecItemCopyMatch - OSStatus: %d", osStatusCode)
	}

	certificateCFData := C.SecCertificateCopyData(secCertificateRef)
	if certificateCFData == 0 {
		return nil, errors.New("GetCertificateDataFromIdentityRef: SecCertificateCopyData: Failed to convert certificate data")
	}
	defer C.CFRelease(C.CFTypeRef(certificateCFData))

	certData := convertCFDataRefToGoBytes(certificateCFData)

	return x509.ParseCertificate(certData)
}

// IdentityWithRefModel ...
type IdentityWithRefModel struct {
	KeychainRef C.CFTypeRef
	Label       string
}

// FindAndValidateIdentity ...
//  IMPORTANT: you have to C.CFRelease the returned items (one-by-one)!!
//             you can use the ReleaseIdentityWithRefList method to do that
func FindAndValidateIdentity(identityLabel string) (*IdentityWithRefModel, error) {
	foundIdentityRefs, err := FindIdentity(identityLabel)
	if err != nil {
		return nil, fmt.Errorf("Failed to find Identity, error: %s", err)
	}
	if len(foundIdentityRefs) < 1 {
		return nil, nil
	}

	// check validity
	var latestIdentityRef *IdentityWithRefModel
	var latestCertificate x509.Certificate

	for _, aIdentityRef := range foundIdentityRefs {
		cert, err := GetCertificateDataFromIdentityRef(aIdentityRef.KeychainRef)
		if err != nil {
			return nil, fmt.Errorf("Failed to read certificate data, error: %s", err)
		}

		if err := certificateutil.CheckValidity(*cert); err != nil {
			log.Warnf("Certificate is not valid, skipping: %s", err)
			continue
		}

		if latestIdentityRef == nil || cert.NotAfter.After(latestCertificate.NotAfter) {
			latestIdentityRef = &aIdentityRef
			latestCertificate = *cert
		}
	}

	return latestIdentityRef, nil
}

// FindIdentity ...
//  IMPORTANT: you have to C.CFRelease the returned items (one-by-one)!!
//             you can use the ReleaseIdentityWithRefList method to do that
func FindIdentity(identityLabel string) ([]IdentityWithRefModel, error) {

	queryDict := C.CFDictionaryCreateMutable(C.kCFAllocatorDefault, 0, nil, nil)
	defer C.CFRelease(C.CFTypeRef(queryDict))
	C.CFDictionaryAddValue(queryDict, unsafe.Pointer(C.kSecClass), unsafe.Pointer(C.kSecClassIdentity))
	C.CFDictionaryAddValue(queryDict, unsafe.Pointer(C.kSecMatchLimit), unsafe.Pointer(C.kSecMatchLimitAll))
	C.CFDictionaryAddValue(queryDict, unsafe.Pointer(C.kSecReturnAttributes), unsafe.Pointer(C.kCFBooleanTrue))
	C.CFDictionaryAddValue(queryDict, unsafe.Pointer(C.kSecReturnRef), unsafe.Pointer(C.kCFBooleanTrue))

	var resultRefs C.CFTypeRef
	osStatusCode := C.SecItemCopyMatching((C.CFDictionaryRef)(queryDict), &resultRefs)
	if osStatusCode != C.errSecSuccess {
		return nil, fmt.Errorf("Failed to call SecItemCopyMatch - OSStatus: %d", osStatusCode)
	}
	defer C.CFRelease(C.CFTypeRef(resultRefs))

	identitiesArrRef := C.CFArrayRef(resultRefs)
	identitiesCount := C.CFArrayGetCount(identitiesArrRef)
	if identitiesCount < 1 {
		return nil, fmt.Errorf("No Identity (certificate + related private key) found in your Keychain")
	}
	log.Debugf("identitiesCount: %d", identitiesCount)

	// filter the identities, by label
	retIdentityRefs := []IdentityWithRefModel{}
	for i := C.CFIndex(0); i < identitiesCount; i++ {
		aIdentityRef := C.CFArrayGetValueAtIndex(identitiesArrRef, i)
		log.Debugf("aIdentityRef: %#v", aIdentityRef)
		aIdentityDictRef := C.CFDictionaryRef(aIdentityRef)
		log.Debugf("aIdentityDictRef: %#v", aIdentityDictRef)

		lablCSting := C.CString("labl")
		defer C.free(unsafe.Pointer(lablCSting))
		vrefCSting := C.CString("v_Ref")
		defer C.free(unsafe.Pointer(vrefCSting))

		labl, err := getCFDictValueUTF8String(aIdentityDictRef, C.CFTypeRef(convertCStringToCFString(lablCSting)))
		if err != nil {
			log.Warnf("FindIdentity: failed to get 'labl' property: %s", err)
			continue
		}
		log.Debugf("labl: %#v", labl)
		if labl != identityLabel {
			continue
		}
		log.Debugf("Found identity with label: %s", labl)

		vrefRef, err := getCFDictValueRef(aIdentityDictRef, C.CFTypeRef(convertCStringToCFString(vrefCSting)))
		if err != nil {
			log.Warnf("FindIdentity: failed to get 'v_Ref' property: %s", err)
			continue
		}
		log.Debugf("vrefRef: %#v", vrefRef)

		// retain the pointer
		vrefRef = C.CFRetain(vrefRef)
		// store it
		retIdentityRefs = append(retIdentityRefs, IdentityWithRefModel{
			KeychainRef: vrefRef,
			Label:       labl,
		})
	}

	return retIdentityRefs, nil
}

//
// --- UTIL METHODS
//

func getCFDictValueRef(dict C.CFDictionaryRef, key C.CFTypeRef) (C.CFTypeRef, error) {
	var retVal C.CFTypeRef
	exist := C.CFDictionaryGetValueIfPresent(dict, unsafe.Pointer(key), (*unsafe.Pointer)(unsafe.Pointer(retVal)))
	// log.Debugf("retVal: %#v", retVal)
	if exist == C.Boolean(0) {
		return 0, errors.New("getCFDictValueRef: Key doesn't exist")
	}
	// return retVal, nil

	return (C.CFTypeRef)(C.CFDictionaryGetValue(dict, unsafe.Pointer(key))), nil
}

func getCFDictValueCFStringRef(dict C.CFDictionaryRef, key C.CFTypeRef) (C.CFStringRef, error) {
	val, err := getCFDictValueRef(dict, key)
	if err != nil {
		return 0, err
	}
	if val == 0 {
		return 0, errors.New("getCFDictValueCFStringRef: Nil value returned")
	}

	if C.CFGetTypeID(val) != C.CFStringGetTypeID() {
		return 0, errors.New("getCFDictValueCFStringRef: value is not a string")
	}

	return C.CFStringRef(val), nil
}

func convertCStringToCFString(cstring *C.char) C.CFStringRef {
	return C.CFStringCreateWithCString(C.kCFAllocatorDefault, cstring, C.kCFStringEncodingUTF8)
}

func getCFDictValueUTF8String(dict C.CFDictionaryRef, key C.CFTypeRef) (string, error) {
	valCFStringRef, err := getCFDictValueCFStringRef(dict, key)
	if err != nil {
		return "", err
	}
	log.Debugf("valCFStringRef: %#v", valCFStringRef)
	if valCFStringRef == 0 {
		return "", errors.New("getCFDictValueUTF8String: Nil value")
	}

	strLen := C.CFStringGetLength(valCFStringRef)
	log.Debugf("strLen: %d", strLen)
	charUTF8Len := C.CFStringGetMaximumSizeForEncoding(strLen, C.kCFStringEncodingUTF8) + 1
	log.Debugf("charUTF8Len: %d", charUTF8Len)
	cstrBytes := make([]byte, charUTF8Len, charUTF8Len)
	if C.Boolean(0) == C.CFStringGetCString(valCFStringRef, (*C.char)(unsafe.Pointer(&cstrBytes[0])), charUTF8Len, C.kCFStringEncodingUTF8) {
		return "", errors.New("getCFDictValueUTF8String: CFStringGetCString: failed to convert value to string")
	}
	return C.GoString((*C.char)(unsafe.Pointer(&cstrBytes[0]))), nil
}
