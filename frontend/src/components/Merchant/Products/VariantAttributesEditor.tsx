"use client";

export type VariantAttributeRow = {
  key: string;
  name: string;
  value: string;
};

const inputClassName =
  "rounded-md border border-gray-3 bg-gray-1 placeholder:text-dark-5 w-full py-2.5 px-5 outline-none duration-200 focus:border-transparent focus:shadow-input focus:ring-2 focus:ring-blue/20";

type VariantAttributesEditorProps = {
  attributes: VariantAttributeRow[];
  onChange: (attributes: VariantAttributeRow[]) => void;
  onAdd: () => void;
};

export default function VariantAttributesEditor({
  attributes,
  onChange,
  onAdd,
}: VariantAttributesEditorProps) {
  const updateRow = (
    index: number,
    field: "name" | "value",
    value: string
  ) => {
    onChange(
      attributes.map((item, i) =>
        i === index ? { ...item, [field]: value } : item
      )
    );
  };

  const removeRow = (index: number) => {
    onChange(attributes.filter((_, i) => i !== index));
  };

  return (
    <section className="mt-8 pt-7.5 border-t border-gray-3">
      <h3 className="font-medium text-lg text-dark mb-5">Attributes</h3>

      <div
        className="hidden sm:grid sm:grid-cols-[1fr_1fr_auto] sm:gap-4 mb-2"
        aria-hidden="true"
      >
        <p className="text-custom-sm text-dark">Attribute name</p>
        <p className="text-custom-sm text-dark">Value</p>
        <span />
      </div>

      <div className="space-y-4">
        {attributes.map((row, index) => (
          <div
            key={row.key}
            className="flex flex-col sm:grid sm:grid-cols-[1fr_1fr_auto] sm:gap-4 sm:items-center gap-2"
          >
            <div>
              <label
                htmlFor={`attr-name-${row.key}`}
                className="block mb-2.5 sm:sr-only"
              >
                Attribute name
              </label>
              <input
                type="text"
                id={`attr-name-${row.key}`}
                placeholder="color"
                value={row.name}
                onChange={(e) => updateRow(index, "name", e.target.value)}
                className={inputClassName}
              />
            </div>
            <div>
              <label
                htmlFor={`attr-value-${row.key}`}
                className="block mb-2.5 sm:sr-only"
              >
                Value
              </label>
              <input
                type="text"
                id={`attr-value-${row.key}`}
                placeholder="Black"
                value={row.value}
                onChange={(e) => updateRow(index, "value", e.target.value)}
                className={inputClassName}
              />
            </div>
            <button
              type="button"
              onClick={() => removeRow(index)}
              className="self-start sm:self-center text-custom-sm font-medium text-red hover:underline py-2.5"
            >
              Remove
            </button>
          </div>
        ))}
      </div>

      <button
        type="button"
        onClick={onAdd}
        className="mt-4 text-custom-sm font-medium text-blue hover:text-blue-dark"
      >
        + Add Attribute
      </button>
    </section>
  );
}
