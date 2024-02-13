export const Login = () => {
  return (
    <div className="flex justify-center items-center flex-col">
      <input
        className="flex w-80 p-1 rounded-md m-2"
        type="text"
        placeholder="Email"
      />
      <input
        className="flex w-80 p-1 rounded-md m-2"
        type="text"
        placeholder="Password"
      />
    </div>
  );
};
