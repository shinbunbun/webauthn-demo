// 登録のサンプル引数
const createCredentialDefaultArgs = {
  publicKey: {
    // Relying Party (a.k.a. - Service):
    rp: {
      name: "bunbun-test-rp"
    },

    // User:
    user: {
      id: new Uint8Array(16),
      name: "john.p.smith@example.com",
      displayName: "John P. Smith"
    },

    pubKeyCredParams: [{
      type: "public-key",
      alg: -7
    }],

    attestation: "direct",

    timeout: 60000,

    challenge: new Uint8Array([ // サーバーから暗号学的にランダムな値が送られていなければならない
      0x8C, 0x0A, 0x26, 0xFF, 0x22, 0x91, 0xC1, 0xE9, 0xB9, 0x4E, 0x2E, 0x17, 0x1A, 0x98, 0x6A, 0x73,
      0x71, 0x9D, 0x43, 0x48, 0xD5, 0xA7, 0x6A, 0x15, 0x7E, 0x38, 0x94, 0x52, 0x77, 0x97, 0x0F, 0xEF
    ]).buffer
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

  // 新しい認証情報の作成/登録
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
    });
}