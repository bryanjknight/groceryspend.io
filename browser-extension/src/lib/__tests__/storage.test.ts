import sinon from "sinon";
import { createAES256LocalStorage } from "../storage";

// create a stub browser
const stubBrowser = {
  storage: {
    local: {
      get: sinon.stub(),
      set: sinon.stub(),
    },
  },
};

describe("encrypted storage tests", () => {
  beforeEach(() => {
    stubBrowser.storage.local.get.reset();
    stubBrowser.storage.local.set.reset();
  });

  test("encrypt value returns same value back", async (done) => {
    const testKey = "test-hash-key";
    const testValue = "token";
    const testCipher = "test-cipher";

    // this was generated using the following commands in node:
    // const cryptoJs = require('crypto-js');
    // cryptoJs.AES.encrypt('token', 'test-cipher').toString()
    const expectedAesEncryptedValue =
      "U2FsdGVkX19Fetm7gnPehSxp/XS9Nm1TL8vAfFmvjUk=";

    // eslint-disable-next-line @typescript-eslint/no-empty-function
    stubBrowser.storage.local.set

      // sinon object matching doesn't do deep check, so we'll
      // have to verify the args passed in ourselves
      .withArgs(sinon.match.object, sinon.match.func)
      .yields();

    stubBrowser.storage.local.get.withArgs([testKey], sinon.match.func).yields({
      [testKey]: expectedAesEncryptedValue,
    });

    // type shenangians to convince typescript that this stub is sufficient
    const encStorage = createAES256LocalStorage(
      (stubBrowser as unknown) as typeof chrome,
      testCipher
    );

    try {
      await encStorage.setEncryptedValue(testKey, testValue);

      sinon.assert.callCount(stubBrowser.storage.local.set, 1);

      // verify the value stored into local storage is not actually the
      // token but the encrypted value. Note encrypted value changes based on the
      // iv, so only verifying we're not storing the unencrypted value
      const actualEncryptedValue = (stubBrowser.storage.local.set.getCall(0)
        .firstArg as Record<string, string>)[testKey];
      expect(actualEncryptedValue !== testValue).toBeTruthy();

      const actualValue = await encStorage.getEncryptedValue(testKey);
      expect(actualValue).toEqual(testValue);
      done();
    } catch (err) {
      done(err);
    }
  });

  test("lookup failure returns null", async (done) => {
    const testKey = "test-hash-key";
    const testCipher = "test-cipher";

    stubBrowser.storage.local.get
      .withArgs([testKey], sinon.match.func)
      .yields({});

    // type shenangians to convince typescript that this stub is sufficient
    const encStorage = createAES256LocalStorage(
      (stubBrowser as unknown) as typeof chrome,
      testCipher
    );

    try {
      const actualValue = await encStorage.getEncryptedValue(testKey);
      expect(actualValue).toBeNull();
      done();
    } catch (err) {
      done(err);
    }
  });
});
