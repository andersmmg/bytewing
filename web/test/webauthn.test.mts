import { describe, it } from "node:test";
import assert from "node:assert/strict";

import {
  credentialToJSON,
  fromBase64Url,
  promptDismissed,
  toBase64Url,
} from "../src/lib/webauthn.ts";

function bytes(...values: number[]): ArrayBuffer {
  const buffer = new ArrayBuffer(values.length);
  new Uint8Array(buffer).set(values);
  return buffer;
}

describe("base64url", () => {
  it("encodes the unpadded URL alphabet WebAuthn identifiers use", () => {
    assert.equal(toBase64Url(bytes(0xfb, 0xff, 0xbe)), "-_--");
    assert.equal(toBase64Url(bytes(0xff, 0xef, 0xff)), "_-__");
    assert.equal(toBase64Url(new ArrayBuffer(0)), "");
    for (const value of [[0], [1, 2, 3], [255, 254, 0, 128, 64]]) {
      const encoded = toBase64Url(bytes(...value));
      assert.equal(encoded, Buffer.from(value).toString("base64url"));
      assert.equal(/^[A-Za-z0-9_-]*$/.test(encoded), true);
    }
  });

  it("round-trips arbitrary bytes", () => {
    const raw = new Uint8Array(300);
    for (let i = 0; i < raw.length; i += 1) raw[i] = (i * 7 + 13) % 256;
    const decoded = new Uint8Array(fromBase64Url(toBase64Url(raw.buffer)));
    assert.deepEqual([...decoded], [...raw]);
  });

  it("accepts padded input, which a server may send", () => {
    const decoded = new Uint8Array(fromBase64Url("AQIDBA=="));
    assert.deepEqual([...decoded], [1, 2, 3, 4]);
  });
});

describe("credentialToJSON", () => {
  it("renders a registration response with its byte fields encoded", () => {
    const credential = {
      id: "AQID",
      rawId: bytes(1, 2, 3),
      type: "public-key",
      authenticatorAttachment: "platform",
      getClientExtensionResults: () => ({}),
      response: {
        clientDataJSON: bytes(10, 20),
        attestationObject: bytes(30, 40, 50),
        getTransports: () => ["internal", "hybrid"],
      },
    } as unknown as PublicKeyCredential;

    const encoded = credentialToJSON(credential) as {
      id: string;
      rawId: string;
      type: string;
      response: Record<string, unknown>;
    };

    assert.equal(encoded.id, "AQID");
    assert.equal(encoded.rawId, "AQID");
    assert.equal(encoded.type, "public-key");
    assert.equal(encoded.response.clientDataJSON, toBase64Url(bytes(10, 20)));
    assert.equal(
      encoded.response.attestationObject,
      toBase64Url(bytes(30, 40, 50)),
    );
    assert.deepEqual(encoded.response.transports, ["internal", "hybrid"]);
  });

  it("renders an assertion response, including the user handle", () => {
    const credential = {
      id: "BAUG",
      rawId: bytes(4, 5, 6),
      type: "public-key",
      authenticatorAttachment: null,
      getClientExtensionResults: () => ({}),
      response: {
        clientDataJSON: bytes(1),
        authenticatorData: bytes(2, 3),
        signature: bytes(4, 5, 6, 7),
        userHandle: bytes(55, 44),
      },
    } as unknown as PublicKeyCredential;

    const encoded = credentialToJSON(credential) as {
      response: Record<string, unknown>;
    };

    assert.equal(encoded.response.authenticatorData, toBase64Url(bytes(2, 3)));
    assert.equal(encoded.response.signature, toBase64Url(bytes(4, 5, 6, 7)));
    assert.equal(encoded.response.userHandle, toBase64Url(bytes(55, 44)));
    assert.equal("attestationObject" in encoded.response, false);
  });
});

describe("promptDismissed", () => {
  it("treats a closed prompt as no answer instead of failure", () => {
    assert.equal(
      promptDismissed(new DOMException("no", "NotAllowedError")),
      true,
    );
    assert.equal(promptDismissed(new DOMException("no", "AbortError")), true);
    assert.equal(
      promptDismissed(new DOMException("no", "NotSupportedError")),
      false,
    );
    assert.equal(promptDismissed(new Error("boom")), false);
  });
});
