import http from 'k6/http';
import { check, sleep } from 'k6';
// vus: 3, duration: '3s',: 91 complete and 0 interrupted : success 100%, avg=96.63ms,p(90)=113.2ms  p(95)=115.14ms
export let options = {
  vus: 3,            // number of virtual users
  duration: '3s',     // total test duration
  thresholds: {
    http_req_duration: ['p(95)<1000'], // 95% of requests must complete < 1000ms
    http_req_failed: ['rate<0.01'],   // less than 1% errors allowed
  },
};

export default function () {
  let res = http.get('https://api.restful-api.dev/objects');
  //if (res.status !== 200) {
  //console.log("Register failed with status: " + res.status);
  //console.log("Body: " + res.body);
//};
  check(res, {
    'status is 200': (r) => r.status === 200,
  });

  // sleep(0.5); // simulate user think-time
}
