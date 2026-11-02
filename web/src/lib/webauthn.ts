/** WebAuthn plumbing for the passkey ceremonies.
 *
 *  A credential is full of ArrayBuffers, which JSON.stringify turns into `{}`,
 *  so every byte field is base64url encoded on the way out and decoded on the
 *  way in — the encoding both the server's parser and the WebAuthn JSON forms
 *  speak. Nothing here decides anything: the server verifies the ceremony,
 *  this file only carries it across the wire.
 */

/** toBase64Url encodes bytes the way WebAuthn identifiers are written:
 *  unpadded, with - and _ instead of + and /. */
export function toBase64Url(buffer: ArrayBuffer): string {
  const bytes = new Uint8Array(buffer);
  let binary = "";
  for (let i = 0; i < bytes.length; i += 1)
    binary += String.fromCharCode(bytes[i]!);
  return btoa(binary)
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=+$/, "");
}

/** fromBase64Url turns such a string back into the bytes the browser API
 *  wants for its BufferSource fields. Padding is optional on the way in. */
export function fromBase64Url(value: string): ArrayBuffer {
  const unpadded = value
    .replace(/=+$/, "")
    .replace(/-/g, "+")
    .replace(/_/g, "/");
  const binary = atob(unpadded + "=".repeat((4 - (unpadded.length % 4)) % 4));
  const buffer = new ArrayBuffer(binary.length);
  const bytes = new Uint8Array(buffer);
  for (let i = 0; i < binary.length; i += 1) bytes[i] = binary.charCodeAt(i);
  return buffer;
}

/** PasskeyCreationOptions is the server's registration options: the same shape
 *  the browser takes, with its byte fields carried as base64url strings. */
export type PasskeyCreationOptions = Omit<
  PublicKeyCredentialCreationOptions,
  "challenge" | "user" | "excludeCredentials"
> & {
  challenge: string;
  user: Omit<PublicKeyCredentialUserEntity, "id"> & { id: string };
  excludeCredentials?: (Omit<PublicKeyCredentialDescriptor, "id"> & {
    id: string;
  })[];
};

/** PasskeyRequestOptions is the same, for a sign-in challenge. */
export type PasskeyRequestOptions = Omit<
  PublicKeyCredentialRequestOptions,
  "challenge" | "allowCredentials"
> & {
  challenge: string;
  allowCredentials?: (Omit<PublicKeyCredentialDescriptor, "id"> & {
    id: string;
  })[];
};

/** toCreationOptions decodes the byte fields so navigator.credentials.create
 *  receives what it expects. */
export function toCreationOptions(
  options: PasskeyCreationOptions,
): PublicKeyCredentialCreationOptions {
  return {
    ...options,
    challenge: fromBase64Url(options.challenge),
    user: { ...options.user, id: fromBase64Url(options.user.id) },
    excludeCredentials: options.excludeCredentials?.map((descriptor) => ({
      ...descriptor,
      id: fromBase64Url(descriptor.id),
    })),
  };
}

/** toRequestOptions is the same for navigator.credentials.get. */
export function toRequestOptions(
  options: PasskeyRequestOptions,
): PublicKeyCredentialRequestOptions {
  return {
    ...options,
    challenge: fromBase64Url(options.challenge),
    allowCredentials: options.allowCredentials?.map((descriptor) => ({
      ...descriptor,
      id: fromBase64Url(descriptor.id),
    })),
  };
}

/** promptDismissed reports whether the browser refused because the person
 *  closed the prompt — not an error, just no answer. */
export function promptDismissed(err: unknown): boolean {
  return (
    err instanceof DOMException &&
    (err.name === "NotAllowedError" || err.name === "AbortError")
  );
}

/** credentialToJSON renders what the browser produced as the plain object the
 *  server parses: every byte field base64url encoded, nothing else touched.
 *
 *  It handles both halves of the protocol — registration (an attestation
 *  object) and sign-in (authenticator data plus a signature) — because the two
 *  differ only in their response fields, and asking which one this is by what
 *  the response carries is more reliable than asking the caller. */
export function credentialToJSON(
  credential: PublicKeyCredential,
): Record<string, unknown> {
  const response = credential.response;

  const body: Record<string, unknown> = {
    clientDataJSON: toBase64Url(response.clientDataJSON),
  };

  if ("attestationObject" in response) {
    // Registration: the browser hands back an attestation object plus the
    // transports the authenticator said it supports.
    const attestation = response as AuthenticatorAttestationResponse;
    body.attestationObject = toBase64Url(attestation.attestationObject);
    // The transports are advisory: they tell the browser which authenticator
    // to offer next time, and are never evidence of anything.
    const transports = attestation.getTransports();
    if (transports.length > 0) body.transports = transports;
  } else {
    // Sign-in: authenticator data and a signature over it.
    const assertion = response as AuthenticatorAssertionResponse;
    body.authenticatorData = toBase64Url(assertion.authenticatorData);
    body.signature = toBase64Url(assertion.signature);
    // Present for a discoverable credential, which is all this flow uses; an
    // absent handle would be refused by the server rather than guessed at.
    if (assertion.userHandle)
      body.userHandle = toBase64Url(assertion.userHandle);
  }

  return {
    // credential.id is already base64url by definition of the WebAuthn API.
    id: credential.id,
    rawId: toBase64Url(credential.rawId),
    type: credential.type,
    authenticatorAttachment: credential.authenticatorAttachment ?? undefined,
    clientExtensionResults: credential.getClientExtensionResults(),
    response: body,
  };
}
