package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	signatureKey := getSignatureJwk(t, signaturePrivateKey)
	require.NotEmpty(t, signatureKey)
	verificationKeys := getVerificationJwks(t)
	require.NotEmpty(t, verificationKeys)

	jwtGenerator, err := NewGenerator(signatureKey, verificationKeys)
	assert.NoError(t, err)
	require.NotEmpty(t, jwtGenerator)
}

func TestGenerator_Sign(t *testing.T) {
	signatureKey := getSignatureJwk(t, signaturePrivateKey)
	require.NotEmpty(t, signatureKey)
	verificationKeys := getVerificationJwks(t)
	require.NotEmpty(t, verificationKeys)

	jwtGenerator, err := NewGenerator(signatureKey, verificationKeys)
	assert.NoError(t, err)
	require.NotEmpty(t, jwtGenerator)

	token := jwt.New()
	err = token.Set(jwt.SubjectKey, subject)
	assert.NoError(t, err)

	signedTokenBytes, err := jwtGenerator.Sign(token)
	assert.NoError(t, err)
	require.NotEmpty(t, signedTokenBytes)
}

func TestGenerator_Verify(t *testing.T) {
	signatureKey := getSignatureJwk(t, signaturePrivateKey)
	require.NotEmpty(t, signatureKey)
	signatureKey2 := getSignatureJwk(t, signaturePrivateKey2)
	require.NotEmpty(t, signatureKey2)
	verificationKeys := getVerificationJwks(t)
	require.NotEmpty(t, verificationKeys)

	tests := []struct {
		Name         string
		SignatureKey jwk.Key
	}{
		{
			Name:         "with signature key 1",
			SignatureKey: signatureKey,
		},
		{
			Name:         "with signature key 2",
			SignatureKey: signatureKey2,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			jwtGenerator, err := NewGenerator(test.SignatureKey, verificationKeys)
			assert.NoError(t, err)
			require.NotEmpty(t, jwtGenerator)

			token := jwt.New()
			err = token.Set(jwt.SubjectKey, subject)
			assert.NoError(t, err)

			signedTokenBytes, err := jwtGenerator.Sign(token)
			assert.NoError(t, err)
			require.NotEmpty(t, signedTokenBytes)

			verifiedToken, err := jwtGenerator.Verify(signedTokenBytes)
			assert.NoError(t, err)
			require.NotEmpty(t, verifiedToken)
			assert.Equal(t, subject, verifiedToken.Subject())
		})
	}
}

func getSignatureJwk(t *testing.T, keyString string) jwk.Key {
	pKey := getPrivateKey(t, keyString)
	key, err := jwk.FromRaw(pKey)
	require.NoError(t, err)
	return key
}

func getPrivateKey(t *testing.T, key string) *rsa.PrivateKey {
	privateKeyBytes, err := base64.RawURLEncoding.DecodeString(key)
	require.NoError(t, err)
	parsedPKey, err := x509.ParsePKCS8PrivateKey(privateKeyBytes)
	require.NoError(t, err)

	pKey, ok := parsedPKey.(*rsa.PrivateKey)
	require.True(t, ok)

	return pKey
}

func getVerificationJwks(t *testing.T) []jwk.Key {
	key1 := getSignatureJwk(t, signaturePrivateKey)
	key2 := getSignatureJwk(t, signaturePrivateKey2)
	return []jwk.Key{key1, key2}
}

var subject = "c21ae0e1-39ad-494f-badd-2d54e072641e"
var signaturePrivateKey = "MIIJQwIBADANBgkqhkiG9w0BAQEFAASCCS0wggkpAgEAAoICAQDXgVHvxGqPUlzVEc9KPufvOSVlmZM9lhEnWMt0XOq5m0wZOvqyxAqmXBOQUzK83MOfV3swwPJM50e0gaG49okNJDUvhYND_JRiA3FcFj85ZsjL9GqwOrc4KpRQVxTcm9w8HdbVmFF_OKYWUt--f9DKj3u21Y7NUMFNF1FzpXDXCGuLxqVknZL_Z4aovpTpHCuTzuASfG-XkvlnTmB5RLCUFagGwFmX1VbS0GNoh-vlgXDaLNjklyG_CQsdVmrr8O6N17W-CdKZJEyXULktd_IRSztPL1U1x1lAAGFpynDNQ1mkclHjN0IZpcjylx5i8rqt3-eJU58yq2HxfkdhozKnwtXWZQ1F1GGlHKAI_6hG1N1_8pIPuiSYy9LNVk_PoUrMxQ9LQbvPeq7yj2Cbso3zbb_h9v_AdBWIaWeX-1Hs3kfqbZhHzpR8WzjDPolNoeu2RXJNIMCyQn3kvBFqiiOSbdbyLIjVOmOUjBtFA217sNaHAcjW9CdsTnClQoktRTaQLNXgeWmD2J4KO2HnqcYLIb57-tea3i6RELlfQUlTFLptK_wXT5NCAbVCnOaq8S4UmZDaGNshUUvA-D7wRFZqTk6JgDMGtC_pU2GcKCO3VDwyGv6zEp83iSRTij_eesEvWVcEEWNz161l_VdXO7fpMPh_3PSyDdvq3hkh46ej8QIDAQABAoICAGMB4b_zEDXKVCX7qa1lmy73pSu5U8EemcDm9Yn_SkN9ioeo5haNJItrj_1li9Di5-jjyxAKBQe51eKjD8anVS25bcnoX_czKoShKkpxWhioFSZGo2FViGmAfmUurMHxxUvFNbcp5H87amqlJnAhzq3RH7hPAu1m5Xfid6RW5LGWB7rOx5ujHS7DxETwUf-K1qZwi9dSXf5YIscIZiAwo6NVE74OTtsHw3zVCmay03i8cDl8EyVqHbHjmLygwDynkyGNcczePGfpGlsGVh0Cly7EznnBuDcd3-4cfqSYwhw7jgqUDvUBpRedZ-Wz8dzpwUQysvAPf_tKa5QEPQ0pahJ6tJhhH-f80S7jqwXK08GhFBlpTdOE4yW8-6jKXV7SdgdACztVTLBnVmcghXrP_BuVoOhJppzRK8lG8qz652ofYMgOzOSxydP0Db_qmcF80NskXSkRFRVzt5i9iC4cqIxfVyks9J0QPt1B9Os2_f7nHslkSqS0lD9v_IsZfMJffHR8dM_WRIYbOYrkcvQz_z-5jSr18sqwFVkb9WxlrdpuHrdrpVV_tv2777DQI1GRtlOZpxOo3OXmjiybJ7Np51zaGnIw_mN_5s4x1EHz5jDGR-H_QW9GbQ9GM40RVs4pj5rHcNabby7ZzuKe0mpEuY7CU3BSjr5osBjl5mU9lwMlAoIBAQD4QjJkPL5pPpqfsnqCPFiAcikpjwabTgF3_xzsX9zj_gFDON4ZCdh2H0B_k7AGA3WModCSF2fCMOyzbtOI4FfhPyZZ7ay3R-BICltszwT5X5BEpjxQjniOICh4craZNApJf4dZg2MfXWBOdzodCaRrFiwCQRTWUh6sGr0H_8efDMQmdCeWISPcKyoO5Z_nUGnLJbEDpneOIAW-Wo7EY3OPac3FM0v0AE-5VOiqQ4qY-i0RkkSLwrCziinBslf8uHEATm_yNQjxTCPUatz1R97lmR-a6M9JcNO2338iTt0Ng6wSFaGza2JOVshNsZ0tGOAywdwMRtIbLym16-LJf7GzAoIBAQDeOaiYzUEsdsRaRdM1OzDsGuszTIG4tbLH3Wacx9CpwQ9OFnAltevLW8jBYJaRsOJHR1KPhc7FJ4bzO3Pfe_APULfnplMr6_bxB4y1Vz8jqO34XDIFkJKMWm6837ld-U8xjLtQ4E6RCepUoUCX-F72QZaVJcAq5uCujcJKXpR4TXTZo6xUPA0-FQH6w7ZXwADhxAX9Ebl5s5MxQeConYiO2noFjXIzYPlH5YkldIbmleR0GfwiNN_DyXiXzlbh4Pd21cfZ2JOJrsZiINPKTZOtbAQ_wshMdMALjSHWtCzCs0b4WyNCdyTAh6mz2xQuinavGuQTnkxxpV6SN3lGp9nLAoIBABfzsw7uuWRIEP0FaEJ2dgd2fDgxP27ueL_OEklP-mzYzeBhdTQvOf4zh7KHWj1KSiYWWpwtu-oFdGDfeXNESdZGlHmqr7ZDLgVlUmrOEmnI6Y9mBn2zMThtK9prHujrF2796d4eCgs1pBwN7sJscruOORLCmrMO2zy5m7FQ4T6cKbSYElWuvtn4JCepyeK0ZHCgI1L51aEVv9gcvpd-DOEyURMMnvBcs1RrN8NtnsqhoIWIeiqNzySTWPICNfEBDo38A1r3-PPm57IP2V-k3oGCY4U7nvwz8Yk8SPTTbQpnwMtB4QcBfkuWnd65GzQFqWPcRlG853qN81VE--1673cCggEBAMJxsRQChQRi52wVrLjnEeeFpkc8qkT0t3oqP57vN6VRSBMLjxVwGOHXbdHGsfjIzTWRMqxiaIoaC_rICpuB1ouQFVqcLipATdKYyIXj0VtidNbb1OkJlzE3761UFN4lRyYT_dLGcfh2tJNYhSx0JqNSwG_AmGTxn6ccYuSv3TlmjNfiXudVpECuIQ1KMkKVvi_NVXAaEjBq8GApRGpFbTeR8zLokQRj1bsTHO2pCGC6xyrPkc5cdW7a2qn54gvCzMUuSbBT0MSoKO2zy504Q_96hD1GMfy0K1XwJ6u1-3RhabfmBvQhTAcqrVKyXvZaMX8GCIsh98F48Ub_Qx6PwAECggEBAOhT9LA4dX_jGxQPBqmJDqa0gEX-RNen4JJPMwgVWpKNKY3uN2EJNdEScQtLbm_mGCHRl712XsfD9kJyygIE9P5X3S6O6CbdtuSKX1md7Vx9ed4tBAghrXBVtv-ZxOJ6Kn1oxRHnNyPbeZlEaAK70IpCJmpzDRINSnSx2WFSdCmh9TmB-jNbRPh32pZUa7_pzmzY81Qg9txLWBRNChI6dlaQYkSCi9aGA1rvkCt8zjvxu7hnyGQmq3Af5mhqhVGcg-QO4XKvT3cbhJQLshZqgTcVsoh03cCTWKRtwNhRZmtOwSULpG6nYuIG3harzh6aWG6NVy6qnJbxSi7tV13yot8"
var signaturePrivateKey2 = "MIIJRAIBADANBgkqhkiG9w0BAQEFAASCCS4wggkqAgEAAoICAQCj9o9Xw8l-y2cbEZ5F8mudA93LpuZlxRMs7oEg0TPhWnwEBtvFCUYCApmKmGBwnLnSCmduc-UQtbPzbMOZ5iDm9BLFUp1uA6YNpTdPTVmcfrK3F9fhcN1qL7tA5wIPry-fLXuCkxQIBF3K4PfzV-R6xj9IkqLWg-dkvF_G2aDbab0ypwho61s35uuPAlc17jsLMcPZ96GhareB9JY14hJu6SAId3tFPFftrLRpMeeXUeicibFszWS_eqXuctnWqxkGeLgzK2Z2g7JxnhTUzIth5PKsD9qImiG6wjxwsA_6jgeag_ZoyPyvEKmHtMXSOquQLlve2FpP3LZa_KeqlNG3UJxXXwssovsEtBQMzrhxEl0Jl_m_FghyT7Y5p7sfpsRP5wxfGdQgKyZ6SlX0L4TXlW9Sg1Zk2yHGOjvDxttmNWALlqR05M9OZVtes3sSmkv4hMJNoV9t_pqxldS3JlD9Dag0K5haMGsqLt-dSpWr9rdgPFpWps0eVJVx46uX9-tR2b_AzIR7wIrQwLnkZUWeQh5ScGarUUU3NEkBJ-zb9EAf9VGq5iWzOwa_LZUo187Wlozf4n-sxlEgBk8eD0AxdQLOiWRpZ61bvA4JOc-_Iyn2J_aHkOJu8-71Etnp7cTjZTjJToOGpX5SOpunD6j26w7jizZtE-wt90tlv0FcBwIDAQABAoICAB9bBnydN4pk8ZnOm7r6qjPDyoWorEToFEuybMVO3KILAM5wVVTv-hBmWOCVVVQT1MFjNfZ8eWDhrsEtmpZy2PXx1SkhLHQehIH1h4dF3o47-IdlKua1A9LLv_6gbtd7BBtnwkftQpZp51nl_eTueQY1pWKGkFd_sB-mmpZXhhiPtxvZr5UI9U_SfFfD9dOddHMmmDK11ZYd52wVzygMzMOjF3onB07tRE9yiKnZXlWk11wgROruuIaZuOPJ2PjhjY6cRXWbpiOh-d0agxdS6pDVDMd03LDk8PBbNt3B_bxHrxKQi-3sCc6c3B7UpkQW7jpLc0xBibveFI1_5byMazJZ9n4aQw-8qQJytOM9a4Bqfy5Umhw3fydF9OfPaJ62oFNhczAFiQDujp_LO6Ct8VwraiOtqJGONptdRR5_QEzoNTm8NTss6kxJ7AY3orq_Shfz_8f7Q07o2npQQ44gsn6VEaIJl1ceGndikC7MvmzOlMsFAN_EQu42W44FG15st_kn19cHXkgaMYuFs_NkMI7FlFAMYRhF26VsZHNa9SwueOza7fIXM6LvJ-5ZndVEHxnjgIX4Qs8M1r68DOqkGc0Ix36ITGlSD7NaBvBywtv6oF1LOXN7oNfcgoKfcddpi6dDhXzu09T_PtLAiXIigfmpHM03loI_hNTGI0Ii9UE5AoIBAQDHErzqrY9qt6Icw9ucee3ivtI1EfKMQMs6orfZupTPPKtC0uyKJv7QOYeptYW72aHVZhAmat4oPH20pagmb7nsq5oJ3ngOVLB5ocLhvTlqXv4o6pDElD6ct7fwDM_afski7P3-EPvCCIGzhte4nihccjVcqdyXxB3125SqIJARds9FqKhhCiJZCIysoUx2Rg9WD2doAzVf4VBZvHQe9g3QhoGrdpx7mTT9kWZ2r2I6ssNis9JWvWqwnmBRBFSJraWoZQilewvgXWbNxjgjVMM2kTce8_nEqB2Ivu0joltSt2zQ9BfkWYfn66iLgBwILoDJOHlT5tbrI3XNr0Fk8eCtAoIBAQDS2ZHxH_KlthiC8Qxxh4n86IeqTldxjAidVbuvjviLmvdPdwr8fEgxscQpOkhOAJayZ9m5bLd7EV8TIRI-VrZs5avrNC2xqVj6ciKOgALGOViaT6Oip7jOQmw-PmhXzGqgemXEGiLcFPfBO7Kyr1ErppZr1bfdghPZwc1cUhzqt05P2McuW1U1EEtTHwvUCXH--JmZrR_OoHBce1t5F6r-V9o4uOdMGV52GkbZ0olwRFWoGst41BA2I5dZtXAjDGfBfmRmnoL1vDoZWmpKHHpaCFXj4W3En2WNNbonIlZYJ-X_9az5zhTRzGMQT61YkVCWVqSgbddG7muBDaTjv-IDAoIBAQCCXD7h3q3f0EiUVZ1mJmIk7ZhxsEMInRV4XD3QkmIII05y__Rts4OBj2rLM1dT3_wd5iwFPE4mQxZ-SUyHfvpdhTHl7IzptYOq4sbfVC5Y_cOpv5D1aa5mqdgFlh42knfcKx1YVn7GXROyIEb7WnZLs25GM-WEbKTB7vy3O4OcLBUnJH6-rQg5DWQxV57ehJpeXM_2SMOW-dDIMqRH8hCx0KLxSUbHmVgeBCz443iLv6w64k6HBprq9YtVAmTpk5C7aQ30b1MjpYZAkeEJIeSDwyw5VMLmiMBuX7iicRskW4Eig_VxTz-0G5nPYCD7KpijLBwnWS675AisnGtSUyIxAoIBAQCNvJHdjENZ3-H0S8O4oYtBxrJD6qvfeWnlEde7-RjdB8wN0BFDjuwc88nQiQxH2x9ySFtUyV9BzHij7ExOOY4h__YkgwvgbN2SZZ0TO7whsjT8bmKqmaKijIYlBWCw_IoE3KKCQ6uBVFsDu3SxpyaieDaPwLf7oFBlxmCdGdm0coqjJC8o216Y8B4ifzE9VSgbZNQkOPuzs6g0kvv3l9Brb3UTQkDBqCAWti1AicW4AUevXGvBCpTnP9-i_1OlS9aHfMZTMWUJeYF4v43JygD5erb6G_TlAt2KIj7DxdJTKmbzPBwORNk7-u_w7A60BeMtXIsICS540RbVRu2756YvAoIBAQCYXS8j_YPMNSQdTpYfzovVIvoAne6t1PuuS9MI3tqx1YN4Y9yDeuM_b1u6cFeKqdb3_MGehhmCiad7giaqu1dfsOmXWJmWoG1WIIXOQyjByWw7WEg1wb00yV10qvWuNeaYWuniPevpKLmwcUX9OqVADw9Di2ik7CjjGzgVglZ3v6cgWMBx4HRkOUS5v4GvlajKAy6hI4il_H_UqQLkk8PbsAdd9BVNs90ArKO9TdSYuANOWrdU9rdf-EPVGvBTy4lIGX1BA36VCOOEzMSGYKlwmuJdaZMq7MggkftQDE63NohWLq4-uJc7QfzPCh4f-GPS0uatzLFrJMGO3pNpMr7a"
