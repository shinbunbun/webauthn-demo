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

const register = () => {

  // メアドと名前を取得
  const email = document.querySelector('#inputEmail').value;
  const displayName = document.querySelector('#inputDisplayName').value;

  



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