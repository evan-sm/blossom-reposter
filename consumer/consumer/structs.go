package main

type Passcode struct {
	Usercode string
	Error    bool
}

var CurrentUsercode Passcode = Passcode{
	Usercode: "",
	Error:    false,
}
