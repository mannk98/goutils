package utils

import "os"

func WindowsIsAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
	/*
		 	if err != nil {
				//        fmt.Println("admin no")
				return false
			}
			//    fmt.Println("admin yes")
			return true
	*/
}
