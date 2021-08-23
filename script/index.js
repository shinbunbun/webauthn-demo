// 登録のサンプル引数
const createCredentialDefaultArgs = {
  publicKey: {

    // RPの情報
    rp: {
      name: '',
      id: 'localhost',
    },

    // ユーザー情報
    user: {
      id: '',
      name: '',
      displayName: '',
    },

    // pubkeyの種類
    // algの種類はhttps://www.iana.org/assignments/cose/cose.xhtml#algorithmsで定義されている
    pubKeyCredParams: [{
      type: 'public-key',
      alg: -7,
    }],

    // RPにAuthenticationデータを渡すかどうか
    attestation: 'direct',

    // タイムアウト時間
    timeout: 60000,

    // サーバーから暗号学的にランダムな値が送られていなければならない
    challenge: '',
  },
};

// ログインのサンプル引数
const getCredentialDefaultArgs = {
  publicKey: {
    timeout: 60000,
    // サーバーから暗号学的にランダムな値が送られていなければならない
    challenge: '',
    rpId: 'localhost',
  },
};

const base64ToArrayBuffer = (base64) => {
  const binary = atob(base64);
  const len = binary.length;
  const bytes = new Uint8Array(len);
  for (let i = 0; i < len; i += 1) {
    bytes[i] = binary.charCodeAt(i);
  }
  return bytes;
};

const base64url2base64 = (base64url) => {
  const base64 = base64url.replace(/-/g, '+').replace(/_/g, '/');
  const padding = base64.length % 4;
  if (padding > 0) {
    return base64 + '===='.slice(padding);
  }
  return base64;
};

// eslint-disable-next-line no-unused-vars
const register = async () => {
  // メアドと名前を取得
  const email = document.querySelector('#inputEmail').value;
  const displayName = document.querySelector('#inputDisplayName').value;

  // eslint-disable-next-line no-undef
  const res = (await axios.post('/register-request', { email, display_name: displayName })).data;
  console.log(res);

  createCredentialDefaultArgs.publicKey.rp.name = res.rp;
  createCredentialDefaultArgs.publicKey.user = {
    id: base64ToArrayBuffer(res.id).buffer,
    name: email,
    displayName,
  };
  createCredentialDefaultArgs.publicKey.challenge = base64ToArrayBuffer(res.challenge).buffer;

  // 新しい認証情報の作成/登録
  navigator.credentials.create(createCredentialDefaultArgs)
    .then(async (cred) => {
      console.log('NEW CREDENTIAL', cred);

      console.log(cred.id);
      console.log(cred.response.clientDataJSON);
      console.log(cred.response.attestationObject);
      console.log(JSON.parse(String.fromCharCode
        .apply(null, new Uint8Array(cred.response.clientDataJSON))));
      // eslint-disable-next-line no-undef
      console.log(cbor.decode(cred.response.attestationObject));

      const authRes = {
        id: base64url2base64(cred.id),
        response: {
          clientDataJSON: JSON.parse(String.fromCharCode
            .apply(null, new Uint8Array(cred.response.clientDataJSON))),
          // eslint-disable-next-line no-undef
          attestationObject: cbor.decode(cred.response.attestationObject),
        },
      };

      authRes.response.attestationObject.attStmt.sig = btoa(
        String.fromCharCode(...authRes.response.attestationObject.attStmt.sig),
      );
      authRes.response.attestationObject.authData = btoa(
        String.fromCharCode(...authRes.response.attestationObject.authData),
      );
      authRes.clientDataJSONString = JSON.stringify(authRes.response.clientDataJSON);

      console.log(authRes);
      console.log(JSON.stringify(authRes));

      // eslint-disable-next-line no-undef
      const regRes = await axios.post('/register', authRes);

      if (regRes.data.verificationStatus === 'succeeded') {
        // eslint-disable-next-line no-restricted-globals
        location.href = '/success-sign-in';
      }

      console.log(btoa(String.fromCharCode(...new Uint8Array(cred.rawId))));
      console.log(base64url2base64(cred.id));
      /* console.log(base64ToArrayBuffer(`${cred.id}==`)); */

      /* const idList = [{
        id: base64ToArrayBuffer(base64url2base64(cred.id)).buffer,
        transports: ['usb', 'nfc', 'ble', 'internal'],
        type: 'public-key',
      }];
      getCredentialDefaultArgs.publicKey.allowCredentials = idList;
      getCredentialDefaultArgs.publicKey.challenge = base64ToArrayBuffer(res.challenge).buffer;
      console.log(getCredentialDefaultArgs);
      return navigator.credentials.get(getCredentialDefaultArgs); */
    })
    .then((assertion) => {
      console.log('ASSERTION', assertion);
    })
    .catch((err) => {
      console.log('ERROR', err);
    });
};

// eslint-disable-next-line no-unused-vars
const login = async () => {
  const email = document.querySelector('#inputEmail').value;
  // eslint-disable-next-line no-undef
  const res = (await axios.post('/login-request', { email, display_name: 'dummy' })).data;
  getCredentialDefaultArgs.publicKey.challenge = base64ToArrayBuffer(res.challenge).buffer;
  const idList = [{
    id: base64ToArrayBuffer(res.id).buffer,
    transports: ['usb', 'nfc', 'ble', 'internal'],
    type: 'public-key',
  }];
  getCredentialDefaultArgs.publicKey.allowCredentials = idList;
  console.log(getCredentialDefaultArgs);
  navigator.credentials.get(getCredentialDefaultArgs)
    .then((resulst) => {
      console.log(resulst);
    })
    .catch((err) => {
      console.log('ERROR', err);
    });
};
