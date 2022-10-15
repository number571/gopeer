#!/bin/bash

str2hex() {
    local str=${1:-""}
    local fmt="%02X"
    local chr
    local -i i
    for i in `seq 0 $((${#str}-1))`; do
        chr=${str:i:1}
        printf "${fmt}" "'${chr}"
    done
}

JSON_DATA='{
        "method":"POST",
        "host":"hidden-echo-service",
        "path":"/echo",
        "head":{
            "Accept": "application/json"
        },
        "body":"aGVsbG8sIHdvcmxkIQ=="
}';

REQUEST_FORMAT="{
        \"receiver\":\"Pub(go-peer/rsa){3082020A0282020100E5B249D1547C5CD9A340CA8DB7DCB789657BAFD39FFA09351DF818CE99CED41451C6D5C7FD87C1CC027C6FC51A8C160DFE5029DA19F02D6E3A237995D0C64FFDEE8423A201663DEDA4574949C76EDE313EE8CC81C7451D602B133936CC045AA73218B2F6359777892A20C4041CAAB279B8718805A56E4C44CEC3BB7B084C1C8EF009E0BAD0F391B4204F96698526CF386B05F0B76CFE6FE278B023C5495D5E728500CFFCF5075DB6AC5214B97D7D14CECAB21E1F79E8AF844999760A040173B0D2DADDC0013A45C3423F127459E0FB70380E8469A7069C31DCD18760E2F0481F1CC2437BEFB801C394FE5D5AC1F985D52D44B980A9E64E46FDA1C3F2A739D2C473057506718EE9903C7E1FD584EE2261593C64543417846B3137098EEBE638AB54149A4EE2839A1DF243B593A3A952A2760885F62986E951C87FEBD99D47FBC7ECF6369DC94663A241ED72BBDB5AD86981ECE2862BE7B32C550A9D9FA141438BA3BA05EF68E52B18F9167CA3D10067CA350D2266385A8470927FAC05A3EBB8D0CB0C9F6430D205F933693478377F892506EAEA2ADE6A9DD736BC54583FF308068DB3B1A37B3ADA22695B27564A4822BDC143381617BCCD0CB7EFB21E1ACA4669B9A72AC3CAC6305407AA16DD50FA8AFEBA680D979242C521A2F9989008FAF6F691A55C34CF0988E4201536F95CF6C92E388278E1F70837C6CCFCC7ACD35272410203010001}\",
        \"hex_data\":\"$(str2hex "$JSON_DATA")\"
}";

curl -i -X PUT -H 'Accept: application/json' http://localhost:7572/do/request --data "${REQUEST_FORMAT}"
