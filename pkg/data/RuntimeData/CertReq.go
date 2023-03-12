package runtimedata

type CertRequest struct {
	Token  string
	CsrStr []byte
}

type CertResponse struct {
	CrtStr []byte
}
