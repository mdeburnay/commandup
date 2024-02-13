export const Login = () => {
  return (
    <div className="flex justify-center flex-col">
      <input
        className="flex w-80 p-1 self-center"
        type="text"
        placeholder="Email"
      />
      <input
        className="flex w-80 p-1 self-center"
        type="text"
        placeholder="Password"
      />
    </div>
  );
};
