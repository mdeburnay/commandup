// Packages
import { useMutation } from "@tanstack/react-query";
import axios from "axios";

// Hooks
import { useState } from "react";

interface ISignupProps {
  email: string;
  password: string;
  username: string;
}

export const Signup = () => {
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [username, setUsername] = useState<string>("");
  const [error, setError] = useState<string>("");
  const [success, setSuccess] = useState<boolean>(false);

  const mutation = useMutation({
    mutationFn: async ({ email, password }: ISignupProps) => {
      return axios.post(
        "http://localhost:8080/api/auth/signup",
        { email, password, username },
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
    mutation.mutate({ email, password, username });
  };

  return (
    <>
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
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <input
          className="flex w-80 p-1 rounded-md m-2"
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <button
          type="submit"
          className="flex p-1 rounded-md m-2 hover:cursor-pointer"
        >
          Signup
        </button>
        {error && <div className="text-red-500">{error}</div>}
        {success && <div>Signup Successful!</div>}
      </form>
    </>
  );
};
