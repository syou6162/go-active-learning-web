import Amplify, {
  Auth,
  API,
} from 'aws-amplify';

Amplify.configure({
  Auth: {
    region: 'us-east-1',
    identityPoolId: 'us-east-1:a61e4d83-cd2a-4bd3-8995-19f9e2726f76',
    userPoolId: 'us-east-1_4Jft7xFAo',
    userPoolWebClientId: '4ia5ifrn456rqg4vr6dqfh7e68',
  },
});

export function signUp(username, password) {
  return Auth.signUp({
    username,
    password,
  });
}

export function signIn(username, password) {
  return Auth.signIn(username, password);
}

export function signOut() {
  return Auth.signOut();
}

export function IsAdmin() {
  return Auth.currentSession().then(_ => true).catch(_ => false);
}
