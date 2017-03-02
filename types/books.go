package types

type Books struct {
	BrandingPackage string `xml:"branding.package,attr"`
	
	BrandingPackageMD5 string `xml:"branding.package.md5,attr"`
	BrandingPackageSHA1 string `xml:"branding.package.sha1,attr"`

	BrandingPackageCompressedSize int `xml:"branding.package.size.compresssed,attr"`
	BrandingPackageRawSize int `xml:"branding.package.size.raw,attr"`
	BrandingPackageTimestamp string `xml:"branding.package.timestamp"`

	Books []Book `xml:"book"`
}
