import { group, check } from "k6";
import http from "k6/http";
import encoding from "k6/encoding";

export default function() {
	group("login", function() {
		let body = JSON.stringify({ username: "haspi", password:"1234" });
		let res = http.post("http://localhost:3440/user/login", body, { headers: { "Content-Type": "application/json" }});
		// Use JSON.parse to deserialize the JSON (instead of using the r.json() method)
		let j = JSON.parse(res.body);
		// Verify response
		check(res, {
			"status is 200": (r) => r.status === 200,
			"is key correct": (r) => j.response === "Success",
		});
	});
	group("market",function(){
		let res = http.get("http://localhost:3440/stock/all");
		check(res, {
			"status is 200": (r) => r.status === 200,
			"check that response is not empty": (r) => r.body.length > 0,
			"check that body size is bigger than 1000 bytes": (r) => r.body.length > 1000,
		});
	});

	group("login-buy-sell process",function(){
		let username = "haspi";
		let password = "1234";
		let body = JSON.stringify({ username: username, password:password });
		let res = http.post("http://localhost:3440/user/login", body, { headers: { "Content-Type": "application/json" }});
		// Use JSON.parse to deserialize the JSON (instead of using the r.json() method)
		let j = JSON.parse(res.body);
		let val_token = j.token;
		body = JSON.stringify({ authToken: val_token,username :username,stockSymbol:"AMZN",operation:"",date:getCurrentDateWithFormat(),expectedStockPrice:"1800.0",amount:1});
		res = http.post("http://localhost:3440/transaction/buy", body, { headers: { "Content-Type": "application/json" }});
		check(res,{
			"was successfull": (r) => r.json().state =="Success",
			"status is 200": (r) => r.status === 200,
		})
		body = JSON.stringify({ authToken: val_token,username :username,stockSymbol:"AMZN",operation:"",date:getCurrentDateWithFormat(),expectedStockPrice:"1800.0", amount:1});
		res = http.post("http://localhost:3440/transaction/sell", body, { headers: { "Content-Type": "application/json" }});
		console.log(res);
		check(res,{
			"was successfull": (r) => r.json().state =="Success",
			"status is 200": (r) => r.status === 200,
		})
	});
}

function decode(token, secret, algorithm) {
    let parts = token.split('.');
    let header = JSON.parse(encoding.b64decode(parts[0], "rawurl"));
    let payload = JSON.parse(encoding.b64decode(parts[1], "rawurl"));
    algorithm = algorithm || algToHash[header.alg];
    if (sign(parts[0] + "." + parts[1], algorithm, secret) != parts[2]) {
        throw Error("JWT signature verification failed");
    }
    return payload;
}

function getCurrentDateWithFormat(){
		var date = new Date();
      var strDate = 'Y-m-d'
        .replace('Y', date.getFullYear())
        .replace('m', date.getMonth() + 1)
        .replace('d', date.getDate());
	return strDate;
}
