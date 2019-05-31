package certificateutil

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

// FilterValidCertificateInfos ...
func FilterValidCertificateInfos(certificateInfos []CertificateInfoModel) []CertificateInfoModel {
	certificateInfosByName := map[string]CertificateInfoModel{}

	for _, certificateInfo := range certificateInfos {
		if certificateInfo.CheckValidity() == nil {
			activeCertificate, ok := certificateInfosByName[certificateInfo.CommonName]
			if !ok || certificateInfo.EndDate.After(activeCertificate.EndDate) {
				certificateInfosByName[certificateInfo.CommonName] = certificateInfo
			}
		}
	}

	validCertificates := []CertificateInfoModel{}
	for _, validCertificate := range certificateInfosByName {
		validCertificates = append(validCertificates, validCertificate)
	}
	return validCertificates
}
