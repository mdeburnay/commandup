import { Link } from "react-router-dom";

interface IButtonProps {
  text: string;
  onClick?: () => void;
  url?: string;
}

export const Button = ({ text, onClick, url }: IButtonProps) => {
  if (url) {
    return (
      <Link
        to={url}
        className="text-blue w-auto cursor-pointer rounded-lg border bg-white px-6 py-1 tracking-wide"
      >
        {text}
      </Link>
    );
  }

  return (
    <button
      className="text-blue w-auto cursor-pointer rounded-lg border bg-white px-6 py-1 tracking-wide"
      onClick={onClick}
    >
      {text}
    </button>
  );
};
