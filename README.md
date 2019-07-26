### VDF (verified delay function) implementation in Golang

This is the VDF used in Harmony Project. It is based on Benjanmin Wesolowski's paper "Efficient verifiable delay functions"(https://eprint.iacr.org/2018/623.pdf).
In this implementation, the VDF function takes 32 bytes as seed and an integer as difficulty.   

Please note that only 2048 integer size for class group variables are supported now.  

The interface of VDF is in src/vdf_go/vdf.go and examples are in src/test/vdf_module_test.go  
