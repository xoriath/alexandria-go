package types

type Books struct {
	BrandingPackage string `xml:"branding.package,attr"`

	BrandingPackageMD5  string `xml:"branding.package.md5,attr"`
	BrandingPackageSHA1 string `xml:"branding.package.sha1,attr"`

	BrandingPackageCompressedSize int    `xml:"branding.package.size.compressed,attr"`
	BrandingPackageRawSize        int    `xml:"branding.package.size.raw,attr"`
	BrandingPackageTimestamp      string `xml:"branding.package.timestamp,attr"`

	Books []Book `xml:"book"`
}

func Unique(in []string) (out []string) {
	m := make(map[string]struct{})

	for _, s := range in {
		m[s] = struct{}{}
	}

	for k := range m {
		out = append(out, k)
	}

	return
}

// Products returns a list containing all defined products
func (b *Books) Products() []string {

	var products []string
	for _, book := range b.Books {
		for _, product := range book.Products {
			products = append(products, product.Name)
		}
	}

	return Unique(products)
}

func (b *Books) Locales(product string) []string {
	var locales []string

	for _, book := range b.Books {
		if book.InProduct(product) {
			locales = append(locales, book.Language)
		}
	}

	return Unique(locales)
}
