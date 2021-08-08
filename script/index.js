// 登録のサンプル引数
const createCredentialDefaultArgs = {
  publicKey: {

    // RPの情報
    rp: {
      name: ''
    },

    // ユーザー情報
    user: {
      id: '',
      name: '',
      displayName: ''
    },

    // pubkeyの種類
    // algの種類はhttps://www.iana.org/assignments/cose/cose.xhtml#algorithmsで定義されている
    pubKeyCredParams: [{
      type: "public-key",
      alg: -7
    }],

    // RPにAuthenticationデータを渡すかどうか
    attestation: "direct",

    // タイムアウト時間
    timeout: 60000,

    // サーバーから暗号学的にランダムな値が送られていなければならない
    challenge: ''
  }
};

// ログインのサンプル引数
const getCredentialDefaultArgs = {
  publicKey: {
    timeout: 60000,
    // サーバーから暗号学的にランダムな値が送られていなければならない
    challenge: ''
  },
};

const base64ToArrayBuffer = (base64) => {
  let binary = atob(base64);
  let len = binary.length;
  let bytes = new Uint8Array(len);
  for (let i = 0; i < len; i++) {
    bytes[i] = binary.charCodeAt(i);
  }
  return bytes;
}

const register = async () => {

  // メアドと名前を取得
  const email = document.querySelector('#inputEmail').value;
  const displayName = document.querySelector('#inputDisplayName').value;

  const res = (await axios.post('/register-request', { email, display_name: displayName })).data;
  console.log(res);

  createCredentialDefaultArgs.publicKey.rp.name = res.rp;
  createCredentialDefaultArgs.publicKey.user = {
    id: base64ToArrayBuffer(res.id).buffer,
    name: email,
    displayName
  };
  createCredentialDefaultArgs.publicKey.challenge = base64ToArrayBuffer(res.challenge).buffer;
  getCredentialDefaultArgs.publicKey.challenge = base64ToArrayBuffer(res.challenge).buffer;

  // 新しい認証情報の作成/登録
  navigator.credentials.create(createCredentialDefaultArgs)
    .then(async (cred) => {
      console.log("NEW CREDENTIAL", cred);

      console.log(cred.id);
      console.log(cred.response.clientDataJSON);
      console.log(cred.response.attestationObject);
      console.log(JSON.parse(String.fromCharCode.apply(null, new Uint8Array(cred.response.clientDataJSON))))
      console.log(cbor.decode(cred.response.attestationObject))

      const authRes = {
        id: cred.id,
        response: {
          clientDataJSON: JSON.parse(String.fromCharCode.apply(null, new Uint8Array(cred.response.clientDataJSON))),
          attestationObject: cbor.decode(cred.response.attestationObject)
        }
      };

      authRes.response.attestationObject.attStmt.sig = btoa(String.fromCharCode(...authRes.response.attestationObject.attStmt.sig));
      authRes.response.attestationObject.authData = btoa(String.fromCharCode(...authRes.response.attestationObject.authData));

      console.log(authRes)
      console.log(JSON.stringify(authRes))

      const res = await axios.post('/register', authRes);

      if (res.data.verificationStatus === "succeeded") {
        location.href = "/success-sign-in"
      }

      /* const idList = [{
        id: cred.rawId,
        transports: ["usb", "nfc", "ble", "internal"],
        type: "public-key"
      }];
      getCredentialDefaultArgs.publicKey.allowCredentials = idList;
      return navigator.credentials.get(getCredentialDefaultArgs); */
    })/* 
    .then((assertion) => {
      console.log("ASSERTION", assertion);
    }) */
    .catch((err) => {
      console.log("ERROR", err);
    });
}