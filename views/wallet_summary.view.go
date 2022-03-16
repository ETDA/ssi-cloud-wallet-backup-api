package views

type WalletSummary struct {
	MyVCCount       int64 `json:"my_vc_count"`
	MyVPCount       int64 `json:"my_vp_count"`
	WaitToSignCount int64 `json:"wait_to_sign_count"`
	IssuedVCCount   int64 `json:"issued_vc_count"`
	IssuedVPCount   int64 `json:"issued_vp_count"`
}
