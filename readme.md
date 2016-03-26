
post /api/rooms -d {name}

put /api/accounts -H auth:token -d {username}

post /api/account {email, password} -> token
post /api/login {email, password} -> token

post /api/clash/{clashid} {code} -> codeid

get /api/code/{codeid} -> score/diff 

post /api/clash/{clash}/code/{code} -> resultid



