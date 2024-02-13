// Components
import { Button } from "./Button";
import { FileUpload } from "./FileUpload";

export const Header = () => {
  return (
    <header className="flex justify-end">
      <FileUpload />
      <Button text="Login" />
    </header>
  );
};
