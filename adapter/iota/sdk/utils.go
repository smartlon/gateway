package sdk


const (
	MWM = 9
	Depth = 3
)

const (
	ENDPOINT ="http://202.117.43.212:14265"
	SEED           = "TX9XRR9SRCOBMTYDTMKNEIJCSZIMEUPWCNLC9DPDZKKAEMEFVSTEVUFTRUZXEHLULEIYJIEOWIC9STAHW"
	sideKeyPublic  = ""
	SIDEKEYPRIVATE = "QOLOACG9BNUYLERQTZPPW9VKIOPDRTPMFZCYWGNVKIZJEYBWJDXASOXNDMZGBNYFVBCFBQBXSCCAFFRIO"
	Message        = "{\"message\":\"Message from Alice\",\"timestamp\":\"2019-4-8 22:41:01\"}"
)

func Must(err error) (bool,error) {
	if err != nil {
		return false,err
	}
	return true,nil
}