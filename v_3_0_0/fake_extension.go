package v_3_0_0

type FakeKeys struct {
	APIServerEncryptionKey []byte
}
type FakeExtension struct {
	Keys FakeKeys
}

func (f *FakeExtension) Files() ([]FileAsset, error) {
	return nil, nil
}

func (f *FakeExtension) Units() ([]UnitAsset, error) {
	return nil, nil
}

func (f *FakeExtension) VerbatimSections() []VerbatimSection {
	return nil
}
