package bridge

import "testing"

func TestClassify(t *testing.T) {
	tests := []struct {
		name string
		auth AuthResult
		want BridgeTrustTier
	}{
		{
			name: "all pass",
			auth: AuthResult{SPF: SPFPass, DKIM: DKIMPass, DMARC: DMARCPass},
			want: TrustTierVerifiedLegacy,
		},
		{
			name: "dkim pass dmarc pass spf fail",
			auth: AuthResult{SPF: SPFFail, DKIM: DKIMPass, DMARC: DMARCPass},
			want: TrustTierVerifiedLegacy,
		},
		{
			name: "dkim fail",
			auth: AuthResult{SPF: SPFPass, DKIM: DKIMFail, DMARC: DMARCPass},
			want: TrustTierSuspicious,
		},
		{
			name: "dmarc fail",
			auth: AuthResult{SPF: SPFPass, DKIM: DKIMPass, DMARC: DMARCFail},
			want: TrustTierSuspicious,
		},
		{
			name: "both fail",
			auth: AuthResult{SPF: SPFFail, DKIM: DKIMFail, DMARC: DMARCFail},
			want: TrustTierSuspicious,
		},
		{
			name: "all none",
			auth: AuthResult{SPF: SPFNone, DKIM: DKIMNone, DMARC: DMARCNone},
			want: TrustTierUnverifiedLegacy,
		},
		{
			name: "dkim pass dmarc none",
			auth: AuthResult{SPF: SPFPass, DKIM: DKIMPass, DMARC: DMARCNone},
			want: TrustTierUnverifiedLegacy,
		},
		{
			name: "dkim none dmarc pass",
			auth: AuthResult{SPF: SPFPass, DKIM: DKIMNone, DMARC: DMARCPass},
			want: TrustTierUnverifiedLegacy,
		},
		{
			name: "spf softfail only",
			auth: AuthResult{SPF: SPFSoftFail, DKIM: DKIMNone, DMARC: DMARCNone},
			want: TrustTierUnverifiedLegacy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Classify(&tt.auth)
			if got != tt.want {
				t.Errorf("Classify() = %d, want %d", got, tt.want)
			}
		})
	}
}
