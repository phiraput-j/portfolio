import http from "k6/http";
import { check, sleep } from "k6";
import { SharedArray } from "k6/data";

export let options = {
  vus: 3,
  duration: "10s",
  thresholds: {
    http_req_failed: ["rate==0"],
    http_req_duration: ["p(95)<1000", "p(99)<2000"],
  },
};

// 🗄️ Pre-generate 100 unique users in memory
const users = new SharedArray("users", function () {
  let arr = [];
  let ts = Date.now(); // same timestamp for this test run
  for (let i = 0; i < 100000; i++) {
    arr.push({
      username: `user_${i}_${ts}@test.com`,
      password: "Pass123!",
    });
  }
  return arr;
});

export default function () {
  let index = (__VU - 1) * 1000 + __ITER;  // unique across VUs
  let user = users[index % users.length];

  // 1️⃣ Register
  let registerRes = http.post(
    "http://localhost:8080/register",
    JSON.stringify(user),
    { headers: { "Content-Type": "application/json" } }
  );

 //console.log("Register response body: " + registerRes.body);
  check(registerRes, {
    "register success": (r) => r.status === 200 || r.status === 201,
    // "valid username in response": (r) =>
    //   r.json("username") && r.json("username") === user.username,
  });
if (registerRes.status !== 200 && registerRes.status !== 201) {
  console.log("Register failed with status: " + registerRes.status);
  console.log("Body: " + registerRes.body);
};

  // 2️⃣ Login
  let loginRes = http.post(
    "http://localhost:8080/login",
    JSON.stringify(user),
    { headers: { "Content-Type": "application/json" } }
  );

  check(loginRes, {
    "login success": (r) => r.status === 200,
    // "token returned": (r) => r.json("token") !== undefined,
  });

  // 3️⃣ Invalid input (error handling)
  // let badRes = http.post(
  //   "http://localhost:8080/register",
  //   JSON.stringify({ username: "", password: "" }),
  //   { headers: { "Content-Type": "application/json" } }
  // );

  // check(badRes, {
  //   "invalid input handled": (r) => r.status === 400 || r.status === 422,
  // });

  sleep(0.5);
}
