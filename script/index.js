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
    // allowCredentials: [newCredential] // 下記参照
    challenge: new Uint8Array([ // サーバーから暗号学的にランダムな値が送られていなければならない
      0x79, 0x50, 0x68, 0x71, 0xDA, 0xEE, 0xEE, 0xB9, 0x94, 0xC3, 0xC2, 0x15, 0x67, 0x65, 0x26, 0x22,
      0xE3, 0xF3, 0xAB, 0x3B, 0x78, 0x2E, 0xD5, 0x6F, 0x81, 0x26, 0xE2, 0xA6, 0x01, 0x7D, 0x74, 0x50
    ]).buffer
  },
};

const register = () => {

  // メアドと名前を取得
  const eMail = document.querySelector('#inputEmail').value;
  const displayName = document.querySelector('#inputDisplayName').value;

  // idを生成
  const array = new Uint32Array(1);
  window.crypto.getRandomValues(array);
  const id = array[0];



  /* // 新しい認証情報の作成/登録
  navigator.credentials.create(createCredentialDefaultArgs)
    .then((cred) => {
      console.log("NEW CREDENTIAL", cred);

      // 通常はサーバーから利用可能なアカウントの認証情報が送られてきますが
      // この例では上からコピーしただけです。
      var idList = [{
        id: cred.rawId,
        transports: ["usb", "nfc", "ble"],
        type: "public-key"
      }];
      getCredentialDefaultArgs.publicKey.allowCredentials = idList;
      return navigator.credentials.get(getCredentialDefaultArgs);
    })
    .then((assertion) => {
      console.log("ASSERTION", assertion);
    })
    .catch((err) => {
      console.log("ERROR", err);
    }); */
}