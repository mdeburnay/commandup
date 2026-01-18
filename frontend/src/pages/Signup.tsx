// Packages
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useSignupMutation } from "../redux/api/apiSlice";

export const Signup = () => {
	const [email, setEmail] = useState<string>("");
	const [password, setPassword] = useState<string>("");
	const [username, setUsername] = useState<string>("");
	const [error, setError] = useState<string>("");

	const navigate = useNavigate();
	const [signup, { isLoading }] = useSignupMutation();

	const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault();
		try {
			await signup({ email, password, username }).unwrap();
			setError("");
			navigate("/login");
		} catch (err: any) {
			setError(err?.data?.message || err.error || "Signup failed");
		}
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
					disabled={isLoading}
					className="flex p-1 rounded-md m-2 hover:cursor-pointer disabled:opacity-50"
				>
					{isLoading ? "Signing up..." : "Signup"}
				</button>
				{error && <div className="text-red-500">{error}</div>}
			</form>
		</>
	);
};
