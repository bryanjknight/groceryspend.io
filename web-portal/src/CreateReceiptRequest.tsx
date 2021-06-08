import React, { useEffect, useRef, useState } from "react";
import { useAuth0 } from "@auth0/auth0-react";
import { ParseReceiptRequest } from "./models";
import { createReceiptRequest } from "./api";

export function CreateReceiptRequest(): JSX.Element {
  const [file, setFile] = useState((null as unknown) as string);
  const [expectedTotalCost, setExpectedTotalCost] = useState(0.0);
  const [submit, setSubmit] = useState(false);
  const { getAccessTokenSilently } = useAuth0();

  const fileInput = useRef((null as unknown) as HTMLInputElement);
  const expectedCostInput = useRef((null as unknown) as HTMLInputElement);

  const getBase64 = (file: File): Promise<string | null> => {
    return new Promise((resolve, reject) => {
      // Make new FileReader
      const reader = new FileReader();

      // Convert the file to base64 text
      reader.readAsDataURL(file);

      // on reader load somthing...
      reader.onload = () => {
        // TODO: better type checking
        if (reader.result === null) {
          resolve(null);
        } else if ((reader.result as ArrayBuffer).byteLength !== undefined) {
          reject(new Error("don't know how to handle arraybuffer"));
        } else {
          const result = reader.result as string;
          resolve(result);
        }
      };
    });
  };

  const handleFileInputChange = async (
    evt: React.ChangeEvent<HTMLInputElement>
  ) => {
    if (!evt.target || !evt.target.files) {
      setFile((null as unknown) as string);
      return;
    }
    const fileBase64 = await getBase64(evt.target.files[0]);
    setFile(fileBase64 || "");
  };

  const handleExpectedCostChange = (
    evt: React.ChangeEvent<HTMLInputElement>
  ) => {
    if (!evt || !evt.target || !evt.target.value) {
      setExpectedTotalCost(0.0);
      return;
    }

    const parsedEventTotalCost = parseFloat(evt.target.value);
    setExpectedTotalCost(parsedEventTotalCost);
  };

  const audience = process.env.REACT_APP_AUDIENCE || "";
  const scope = "read:users";

  useEffect(() => {
    (async () => {
      if (submit) {
        try {
          // get the bearer token
          const accessToken = await getAccessTokenSilently({
            audience,
            scope,
            timeoutInSeconds: 60 * 60,
          });

          const req: ParseReceiptRequest = {
            data: file,
            parseType: 2, // image = 2
            timestamp: new Date(),
            expectedTotal: expectedTotalCost,
          };
          await createReceiptRequest(req)(accessToken);
          setSubmit(false);
          setFile((null as unknown) as string);
          setExpectedTotalCost((null as unknown) as number);

          if (fileInput) {
            fileInput.current.value = "";
          }
          if (expectedCostInput) {
            expectedCostInput.current.value = "";
          }
        } catch (error) {
          console.error(error);
        }
      }
    })();
  }, [audience, expectedTotalCost, file, getAccessTokenSilently, submit]);

  return (
    <div>
      <div>
        <input
          ref={fileInput}
          type="file"
          name="file"
          onChange={handleFileInputChange}
        />{" "}
      </div>
      <div>
        <label>Enter total cost (shown at bottom of receipt)</label>
        <input
          type="text"
          name="expectedCost"
          ref={expectedCostInput}
          onChange={handleExpectedCostChange}
        />
      </div>
      <div>
        <button type="submit" name="submit" onClick={() => setSubmit(true)}>
          Submit
        </button>
      </div>
    </div>
  );
}
