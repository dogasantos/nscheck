package nscheck

import (
	"fmt"
	"sync"

	"github.com/miekg/dns"
	"github.com/projectdiscovery/retryabledns"
	"github.com/xgfone/netaddr"
)

var defaultResolvers = []string{
	"1.1.1.1:53", // Cloudflare
	"1.0.0.1:53", // Cloudflare
	"8.8.8.8:53", // Google
	"8.8.4.4:53", // Google
}

/*
Line 20: warning: exported function CheckIpv4Type should have comment or be unexported (golint)
now you have a comment.
*/
func CheckIpv4Type(ipaddr string) string {
	var result string
	ipv4 := netaddr.MustNewIPAddress(ipaddr)
		if ipv4.IsIPv4() {
			if ipv4.IsLoopback() {
				result = fmt.Sprintf("%s:loopback",ipaddr)
			} else {
				if ipv4.IsReserved() {
					result = fmt.Sprintf("%s:reserved",ipaddr)
				} else {

					if ipv4.IsPrivate() {
						result = fmt.Sprintf("%s:private",ipaddr)
					} else {
						result = fmt.Sprintf("%s:public",ipaddr)
					}
				}
			}
			
		} else {
			result = fmt.Sprintf("%s:invalid",ipaddr)
		}
	
	return result
}
func doResolve(hostname string, resolvers []string){
	var responseValue []string

	dnsClient := retryabledns.New(resolvers, 2)
	dnsResponse, _ := dnsClient.Query(hostname, dns.TypeA)
	responseValue = dnsResponse.A
	return responseValue
}

func CheckResolvers(resolvers []string, TrustedNs string, wg * sync.WaitGroup, verbose bool) {
	var iptype string
	var tr []string
	defer wg.Done()
	
	// 1) verificar trusted host com trusted server, pegar o ip
	if len(TrustedNs) >0 {
		tr[0] = TrustedNs
	} else {
		tr = defaultResolvers
	}
	rdata := doResolve("nscheck.pingback.me",tr)


	/*
	for _, ipaddr := range resolvers {
		iptype = CheckIpv4Type(ipaddr)
		if strings.Contains("public",strings.Split(iptype,":")[1]) {
			if strings.Split(ipaddr, ".")[0] == strings.Split(prefix, ".")[0] {
				if verbose == true{
					fmt.Printf("  + Checking: %s and %s\n",ipaddr,prefix)
				}
				// testa se o ip eh valido/invalido/publico/
				if len(ipaddr) > 4 {
					result:=CheckIp(ipaddr)
					if strings.Contains("public",strings.Split(result, ":")[1]) {
						iptest := net.ParseIP(ipaddr)
						if cidrAddr.Contains(iptest) == true {
							fmt.Println(ipaddr)
						}
					}
				}
			}
		}
	}
	*/
	fmt.Println(rdata)

}
