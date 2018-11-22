package dispatch_test

// This includes the tests for dispatch. The important part that needs to be checked is the online finder, since it has the most possible paths.

/*
// Without exclusions, test
addresses, err := dispatch.GetOnlineAddresses(1, []api.Address{})
if err != nil {
  logging.Log(1, err)
}
fmt.Printf("%#v\n", addresses)

// With exclusions, test
var excludedAddr api.Address
excludedAddr.Location = "127.0.0.1"
excludedAddr.Sublocation = ""
excludedAddr.Port = 8000
exclusions := []api.Address{excludedAddr}
addresses, err2 := dispatch.GetOnlineAddresses(1, exclusions)
if err2 != nil {
  logging.Log(err2)
}
fmt.Printf("%#v\n", addresses)


*/
