package certificateutil

import "sort"

// FilterCertificateInfoModelsByFilterFunc ...
func FilterCertificateInfoModelsByFilterFunc(certificates []CertificateInfoModel, filterFunc func(certificate CertificateInfoModel) bool) []CertificateInfoModel {
	filteredCertificates := []CertificateInfoModel{}

	for _, certificate := range certificates {
		if filterFunc(certificate) {
			filteredCertificates = append(filteredCertificates, certificate)
		}
	}

	return filteredCertificates
}

// ValidCertificateInfo contains the certificate infos filtered as valid, invalid and duplicated common name certificates
type ValidCertificateInfo struct {
	ValidCertificates,
	InvalidCertificates,
	DuplicatedCertificates []CertificateInfoModel
}

// FilterValidCertificateInfos filters out invalid and duplicated common name certificaates
func FilterValidCertificateInfos(certificateInfos []CertificateInfoModel) ValidCertificateInfo {
	var invalidCertificates []CertificateInfoModel
	nameToCerts := map[string][]CertificateInfoModel{}
	for _, certificateInfo := range certificateInfos {
		if certificateInfo.CheckValidity() != nil {
			invalidCertificates = append(invalidCertificates, certificateInfo)
			continue
		}

		nameToCerts[certificateInfo.CommonName] = append(nameToCerts[certificateInfo.CommonName], certificateInfo)
	}

	var validCertificates, duplicatedCertificates []CertificateInfoModel
	for _, certs := range nameToCerts {
		if len(certs) == 0 {
			continue
		}

		sort.Slice(certs, func(i, j int) bool {
			return certs[i].EndDate.Before(certs[j].EndDate)
		})
		validCertificates = append(validCertificates, certs[0])
		if len(certs) > 1 {
			duplicatedCertificates = append(duplicatedCertificates, certs[1:]...)
		}
	}

	return ValidCertificateInfo{
		ValidCertificates:      validCertificates,
		InvalidCertificates:    invalidCertificates,
		DuplicatedCertificates: duplicatedCertificates,
	}
}
