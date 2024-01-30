// Dependencies
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

// Styles
import "./index.css";

// Components
import { Container } from "./components/Container";

// Pages
import { Home } from "./pages/Home";

// Create a client
const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Container>
        <Home />
      </Container>
    </QueryClientProvider>
  );
}

export default App;
