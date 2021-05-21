import { useApi } from "./use-api";
import React, { useEffect, useState } from "react";
import { Loading } from "./Loading";
import { Error } from "./Error";
import { RouteComponentProps } from "react-router-dom";
import { ReceiptItem, ReceiptDetail, Category } from "./models";
import { getAllCategories, getReceiptDetails, patchItemCategory } from "./api";
import {
  EditableCell,
  EditableCellProps,
} from "./components/tables/EditableCell";
import { Dropdown, DropdownProps } from "./components/forms/Dropdown";
import { useAuth0 } from "@auth0/auth0-react";

export const ReceiptItemRow = (props: {
  receiptID: string;
  item: ReceiptItem;
  catData: Category[];
}) => {
  const [category, setCategory] = useState(props.item.Category);
  const [stale, setStale] = useState(false);
  const { getAccessTokenSilently } = useAuth0();
  const { receiptID, item } = props;

  // eslint-disable-next-line react-hooks/exhaustive-deps
  const dummyCategory = { ID: 0, Name: "unknown" };

  const audience = process.env.REACT_APP_AUDIENCE || "";
  const scope = "read:users";

  // TODO: how can I refactor this so that it's more reuseable
  useEffect(() => {
    (async () => {
      if (stale) {
        try {
          // get the bearer token
          const accessToken = await getAccessTokenSilently({ audience, scope, timeoutInSeconds: 60*60 });
          await patchItemCategory(receiptID, item, category || dummyCategory)(accessToken);
          setStale(false);
        }
        catch (error) {
          console.error(error);
        }
      }
    })();
    
  }, [audience, category, dummyCategory, getAccessTokenSilently, item, receiptID, stale])

  const handleCategoryEdit = (item: ReceiptItem, currentCategory: Category) => {
    setCategory(currentCategory);
    setStale(true);

    // call API to update
    console.log(
      `Updating item ${item.ID} to have category ${currentCategory.Name}`
    );
  };

  const getEditableCellProps = (
    item: ReceiptItem
  ): EditableCellProps<Category> => {
    return {
      id: `${item.ID}-cat-editor`,
      value: item.Category || { ID: 0, Name: "unknown" },
      className: "",
      editorFactory: (handleChange, handleOnBlur) => {
        // create dropdown props
        const dropdownProps: DropdownProps<Category> = {
          id: `${item.ID}-dropdown`,
          mapOptionsToSelectItems: (c: Category) => ({
            label: c.Name,
            value: c.ID.toString(),
          }),
          onSelect: (c: Category) => handleChange(c),
          onBlur: () => handleOnBlur(),
          options: props.catData,
        };
        return <Dropdown {...dropdownProps} />;
      },
      onValueChange: (c: Category) => {
        handleCategoryEdit(item, c);
      },
      valueLabelMaker: (c: Category) => c.Name,
    };
  };

  return (
    <tr key={item.ID}>
      <td>{item.Name}</td>
      <EditableCell {...getEditableCellProps(item)} />
      <td>${item.TotalCost.toFixed(2)}</td>
    </tr>
  );
};

export function ReceiptDetails(props: RouteComponentProps): JSX.Element {
  const params = props.match.params;
  const receiptID = "ID" in params ? params["ID"] : "";
  const { loading, error, data } = useApi<ReceiptDetail | null>(
    getReceiptDetails({ receiptUuid: receiptID }),
    {
      audience: process.env.REACT_APP_AUDIENCE || "",
      scope: "read:users",
      // mode: "cors",
      // credentials: "include",
    }
  );

  // TODO: memoize this call
  const { loading: catLoading, error: catError, data: catData } = useApi<
    Category[]
  >(getAllCategories(), {
    audience: process.env.REACT_APP_AUDIENCE || "",
    scope: "read:users",
    // mode: "cors",
    // credentials: "include",
  });

  if (loading || catLoading) {
    return <Loading />;
  }

  if (error) {
    return <Error message={error.message} />;
  }
  if (catError) {
    return <Error message={catError.message} />;
  }

  if (!data) {
    return <div>Receipt not found</div>;
  }

  return (
    <table className="table">
      <thead>
        <tr>
          <th scope="col">Item</th>
          <th scope="col">Category</th>
          <th scope="col">Cost</th>
        </tr>
      </thead>
      <tbody>
        {data.Items?.map((item: ReceiptItem, i: number) => (
          <ReceiptItemRow
            key={item.ID}
            receiptID={receiptID}
            item={item}
            catData={catData || []}
          />
        ))}
      </tbody>
    </table>
  );
}
