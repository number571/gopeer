# gopeer

> Framework for create secure decentralized applications. Version: 1.3

## Research Article
* Theory of the structure of hidden systems: [hiddensystems.pdf](https://github.com/Number571/gopeer/blob/master/hiddensystems.pdf "TSHS");

## Framework based applications
* Hidden Lake: [github.com/number571/HiddenLake](https://github.com/number571/HiddenLake "HL");
* Hidden Lake Service: [github.com/number571/gopeer/HLS](https://github.com/number571/gopeer/tree/master/cmd/HLS "HLS");
* Hidden Email Service: [github.com/number571/gopeer/HES](https://github.com/number571/gopeer/tree/master/cmd/HES "HES");

## Description
> Part from "Theory of the structure of hidden systems" (Translated) [pages 8,9]

If we assume that there are only three nodes `{A, B, C}` in the network (where one of them is the sender - A) and the network itself is based on the sixth stage of anonymity without information polymorphism, then in this case and under this condition it is extremely problematic to determine the true recipient until he gives himself out as a response to the request (since the response will be a completely new packet, different from all the others). Now, if we assume that there is a possibility of information polymorphism, that is, the probability of its routing, then the stage of merging the properties of receiving and sending begins, forming an anticipation. So, for example, if polymorphism exists, then there will be three stages: `(A → B OR A → C) → (B → C OR C → B) → (B → A OR C → A)`, but if polymorphism does not exist, then there will be two stages: `(A → B OR A → C) → (B → A OR C → A)`. It is assumed that the system knows only the sender of the information (initiator), while the recipient is not defined. It follows that if polymorphism is a static value (that is, it will always exist or not exist at all), then determining the recipient will be an easy task (provided that it always responds to the initiator). But, if polymorphism has a probabilistic value, then the line between sending and receiving will be erased, merged, inverted, which will lead to different interpretations of the analyzed actions: `request (1) - response (1) - request (2)` or `request (1) - routing (1) - response (1)`. But in this case, the property of `hyperthelia` (over the end) arises, where request (2) does not receive its answer (2), which again leads to the possibility of deterministic determination of subjects. Now, if we align the number of polymorphism actions (the number of packet routing) k and the number of actions without it n (which is always a constant n = 2), in other words, adhere to the formula `GCD (k, 2) = 2` (where GCD is the greatest common divisor) , then we get the maximum uncertainty, aleatoryness at a constant k = 2, which can be reduced to the following minimum set of polymorphism actions: `(A → B OR A → C) → (B → C OR C → B) → (B → C OR C → B) → (B → A OR C → A)`. As a result, all actions can be interpreted as two completely self-sufficient processes: `request (1) - response (1) - request (2) - response (2)` or `request (1) - routing (1) - routing (~ 1) - response (1)`, which in turn leads to the uncertainty of sending and receiving information from the traffic analysis of the entire network. And so `response (1) = routing (1)`, `request (2) = routing (~ 1)`, and also `response (2) = response (1) = routing (2)`, where the last incremental routing (2) comes from request (2). The problem, in this case, is only the request (1), created by the initiator of the connection, which will always be interpreted deterministically. But here it is worth noting that with subsequent actions, this problem will always fade away due to the increasing entropy, leading to chaotic actions. For example, at the next step, there will be an ambiguity of the form `request (3) = request (2) = routing (~ 2)`, which means that the sender is not uniquely identified. 

## Template
> Creating a node with a port setting to accept data and a listening function 
```go
package main

import (
	"fmt"

	gp "github.com/number571/gopeer"
	cr "github.com/number571/gopeer/crypto"
	lc "github.com/number571/gopeer/local"
	nt "github.com/number571/gopeer/network"
)

func init() {
	gp.Set(gp.SettingsType{
		"AKEY_SIZE": uint(1 << 10),
		"SKEY_SIZE": uint(1 << 4),
	})
}

func main() {
	fmt.Println("Node is listening...")
	client := lc.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint)))
	nt.NewNode(client).
		Handle([]byte("/msg"), msgRoute).Listen(":8080")
	// ...
}

func msgRoute(client *lc.Client, msg *lc.Message) []byte {
	hash := cr.LoadPubKey(msg.Head.Sender).Address()
	fmt.Printf("[%s] => '%s'\n", hash, msg.Body.Data)
	return msg.Body.Data
}
```

## Entropy increase
> Example of how the sixth stage of anonymity generates probabilistic polymorphism

1. request(1)[`I, II`] <br>
	(_A → B OR A → C_) <br>
2. routing(1)[`I`] = response(1)[`II`] <br>
	(_B → C OR C → B_) <br>
	_OR_ <br>
	(_B → A OR C → A_) <br>
3. routing(~1)[`I`] = request(2)[`II, III`] <br>
	(_B → C OR C → B_) <br>
	_OR_ <br>
	( (_B → A OR B → C_) _OR_ (_C → A OR C → B_) ) <br>
4. response(1)[`I`] = response(2)[`II`] = routing(2)[`III`] <br>
	(_B → A OR C → A_) <br>
	_OR_ <br>
	( (_A → B OR C → B_) _OR_ (_A → C OR B → C_) ) <br>
	_OR_ <br>
	( (_A → C OR C → A_) _OR_ (_A → B OR B → A_) ) <br>
5. ... 

## Need to do
> Pages from “Theory of the structure of hidden systems” 

At the moment, the framework is able to recreate the fourth stage of anonymity, but is not suitable for the sixth. This is due to the three pitfalls of the sixth stage of anonymity that need to be corrected. The list is as follows:
1. Request time. You need to implement a simulation of packet generation time, either on a request-based or routing-response basis [page 10].
2. The period of states. This problem should be solved dynamically by the user and based on the framework it is quite impossible to fix the vulnerability [page 10].
3. Package size. It is necessary to adjust all packets to the dynamic (random) value when sending [pages 10,11,12,13]. 
