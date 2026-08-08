const mockStats = [
  { label: "Revenue", value: "€12,450.00", change: "+12.5%" },
  { label: "Orders", value: "156", change: "+8.2%" },
  { label: "Products", value: "23", change: "0%" },
  { label: "Available Balance", value: "€2,340.00", change: null },
] as const;

const mockOrders = [
  {
    id: "ORD-001234",
    customer: "John Smith",
    date: "2026-08-08",
    status: "Completed",
    total: "€89.99",
  },
  {
    id: "ORD-001235",
    customer: "Sarah Johnson",
    date: "2026-08-08",
    status: "Processing",
    total: "€124.50",
  },
  {
    id: "ORD-001236",
    customer: "Mike Davis",
    date: "2026-08-07",
    status: "Pending",
    total: "€45.00",
  },
  {
    id: "ORD-001237",
    customer: "Emma Wilson",
    date: "2026-08-07",
    status: "Completed",
    total: "€199.99",
  },
  {
    id: "ORD-001238",
    customer: "Chris Brown",
    date: "2026-08-06",
    status: "Processing",
    total: "€67.50",
  },
] as const;

function statusClass(status: string) {
  switch (status) {
    case "Completed":
      return "bg-green-light-6 text-green";
    case "Processing":
      return "bg-blue-light-5 text-blue";
    case "Pending":
      return "bg-yellow-light-4 text-yellow";
    default:
      return "bg-gray-1 text-dark-2";
  }
}

const OverviewTab = () => {
  return (
    <div className="w-full min-w-0 bg-white rounded-xl shadow-1 py-9.5 px-4 sm:px-7.5 xl:px-10">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-7.5">
        <h2 className="font-medium text-xl sm:text-2xl text-dark">Overview</h2>
        <button
          type="button"
          className="inline-flex items-center justify-center font-medium text-white bg-blue py-3 px-7 rounded-md ease-out duration-200 hover:bg-blue-dark"
        >
          + Add Product
        </button>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4 sm:gap-5 mb-7.5">
        {mockStats.map((stat) => (
          <div
            key={stat.label}
            className="rounded-xl border border-gray-3 bg-gray-1 p-5"
          >
            <p className="text-custom-sm text-dark-4 mb-1.5">{stat.label}</p>
            <p className="font-medium text-xl text-dark mb-1">{stat.value}</p>
            {stat.change !== null && (
              <p
                className={`text-custom-sm ${
                  stat.change.startsWith("+")
                    ? "text-green"
                    : "text-dark-4"
                }`}
              >
                {stat.change}
              </p>
            )}
          </div>
        ))}
      </div>

      <div>
        <div className="flex items-center justify-between gap-3 mb-5">
          <h3 className="font-medium text-lg text-dark">Recent Orders</h3>
          <button
            type="button"
            className="text-custom-sm font-medium text-blue hover:underline"
          >
            View All Orders
          </button>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full min-w-[640px] text-left">
            <thead>
              <tr className="border-b border-gray-3">
                <th className="pb-4 pr-4 font-medium text-custom-sm text-dark-4">
                  Order ID
                </th>
                <th className="pb-4 pr-4 font-medium text-custom-sm text-dark-4">
                  Customer
                </th>
                <th className="pb-4 pr-4 font-medium text-custom-sm text-dark-4">
                  Date
                </th>
                <th className="pb-4 pr-4 font-medium text-custom-sm text-dark-4">
                  Status
                </th>
                <th className="pb-4 font-medium text-custom-sm text-dark-4 text-right">
                  Total
                </th>
              </tr>
            </thead>
            <tbody>
              {mockOrders.map((order) => (
                <tr key={order.id} className="border-b border-gray-3 last:border-0">
                  <td className="py-4 pr-4 text-custom-sm font-medium text-dark">
                    {order.id}
                  </td>
                  <td className="py-4 pr-4 text-custom-sm text-dark">
                    {order.customer}
                  </td>
                  <td className="py-4 pr-4 text-custom-sm text-dark-4">
                    {order.date}
                  </td>
                  <td className="py-4 pr-4">
                    <span
                      className={`inline-flex rounded-full px-2.5 py-1 text-2xs font-medium ${statusClass(
                        order.status
                      )}`}
                    >
                      {order.status}
                    </span>
                  </td>
                  <td className="py-4 text-custom-sm font-medium text-dark text-right">
                    {order.total}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default OverviewTab;
