import React from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router";
import WishlistPage from "./pages/WishlistPage";
import WishDetailPage from "./pages/WishDetailPage";
import ListsPage from "./pages/ListsPage";

const App: React.FC = () => {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<ListsPage />} />
                <Route path="/wishlist/new" element={<WishlistPage />} />
                <Route path="/wishlist/:id" element={<WishlistPage />} />
                <Route path="/wishlist/:id/wish/new" element={<WishDetailPage />} />
                <Route path="/wishlist/:id/wish/:wishid" element={<WishDetailPage />} />
                <Route path="*" element={<Navigate to="/" />} />
            </Routes>
        </BrowserRouter>
    );
};

export default App;