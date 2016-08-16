package endpoint

func (e *endpoint) sendScanOptions() (string, error) {
	// TODO(bga): Generate options dynamically based on
	// SANE's scanner options.
	request := `
<?xml version="1.0" encoding="UTF-8" ?>
<root>
    <S2PC_AppList>
        <List>
            <AppIndex Value="1" />
            <AppName Value="JPEG-Colorido-300" />
            <AppType Value="MAC" />
            <Resolution Value="DPI_300" />
            <Color Value="COLOR_TRUE" />
            <FileFormat Value="FORMAT_JPEG" />
            <ScanSize Value="SIZE_A4" />
            <DuplexScan Value="DUPLEX_OFF" />
            <Orientation Value="ORIENTATION_SIDEWAY" />
        </List>
    </S2PC_AppList>
</root>`

	return formUpload(e.s.IP(), "/IDS/ScanFaxToPC.cgi", "scantopc", request,
		true)
}

func (e *endpoint) getUserScanOptions() (string, error) {
	return formUpload(e.s.IP(), "/IDS/UserSelect.xml", "scantopc",
		"", false)
}
