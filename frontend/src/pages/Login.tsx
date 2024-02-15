// Packages
import { useMutation } from "@tanstack/react-query";
import axios from "axios";

// Hooks
import { useState } from "react";

interface ILoginProps {
  email: string;
  password: string;
}

export const Login = () => {
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [error, setError] = useState<string>("");
  const [success, setSuccess] = useState<boolean>(false);

  const mutation = useMutation({
    mutationFn: async ({ email, password }: ILoginProps) => {
      return axios.post(
        "http://localhost:8080/api/auth/login",
        { email, password },
        {
          headers: {
            "Content-Type": "application/json",
          },
        }
      );
    },
    onError: (error) => {
      setError(error.message);
    },
    onSuccess: () => {
      setSuccess(true);
      setError("");
    },
  });

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    mutation.mutate({ email, password });
  };

  return (
    <form
      className="flex justify-center items-center flex-col"
      onSubmit={handleSubmit}
    >
      <input
        className="flex w-80 p-1 rounded-md m-2"
        type="text"
        placeholder="Email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />
      <input
        className="flex w-80 p-1 rounded-md m-2"
        type="text"
        placeholder="Password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      <button
        type="submit"
        className="flex p-1 rounded-md m-2 hover:cursor-pointer"
      >
        Login
      </button>
      {error && <div className="text-red-500">{error}</div>}
      {success && <div>Login Successful!</div>}
    </form>
  );
};
