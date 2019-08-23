package main

func main() {
	a := App{}
	a.Initialize(getRedisEnv())
	a.Run(getAppEnv().Addr)
}
