import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 10,
  duration: '30s',
};

export default function () {
  const payload = JSON.stringify({
    clusterUID: Math.random().toString(16).slice(2),
    bootstrapToken: 'poc-token',
    metadata: {
      region: 'kind-west',
      tenant: 'demo',
    },
  });

  const res = http.post('http://localhost:8080/api/v1/register', payload, {
    headers: { 'Content-Type': 'application/json' },
  });

  check(res, {
    'status was 200': (r) => r.status === 200,
  });

  sleep(1);
}
