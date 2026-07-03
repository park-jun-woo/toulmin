//go:build fakecodec

//ff:func feature=cli type=command control=sequence
//ff:what fakeEncodeCodec.Locale — always report locale "ko" (the fixture MD's detected locale)
package tanglcli

// Locale always reports "ko", matching the locale detectLocale recovers from the
// `## tangl:주제` fixture header used by TestRunEncodeCodec. It is "ko" (not "en")
// so the fake does not collide with the real en codec registerCodecs installs.
func (f *fakeEncodeCodec) Locale() string { return "ko" }
