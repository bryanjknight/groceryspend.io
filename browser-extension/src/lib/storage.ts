import CryptoJS from "crypto-js";

export interface EncryptedStorage {
  getEncryptedValue(key: string): Promise<string | null>;
  setEncryptedValue(key: string, value: string): Promise<void>;
}

const encryptAES256 = (browser: typeof chrome, cipher: string) => (
  key: string,
  value: string
): Promise<void> => {
  return new Promise<void>((resolve) => {
    browser.storage.local.set(
      {
        [key]: CryptoJS.AES.encrypt(value, cipher).toString(),
      },
      () => {
        resolve();
      }
    );
  });
};

const decryptAES256 = (browser: typeof chrome, cipher: string) => (
  key: string
): Promise<string | null> => {
  return new Promise<string | null>((resolve) => {
    browser.storage.local.get([key], (items) => {
      if (!(key in items)) {
        resolve(null);
      } else {
        const cipherText = items[key];
        resolve(
          CryptoJS.AES.decrypt(cipherText, cipher).toString(CryptoJS.enc.Utf8)
        );
      }
    });
  });
};

export const createAES256LocalStorage = (
  browser: typeof chrome,
  cipher: string
): EncryptedStorage => {
  return {
    getEncryptedValue: decryptAES256(browser, cipher),
    setEncryptedValue: encryptAES256(browser, cipher),
  };
};
