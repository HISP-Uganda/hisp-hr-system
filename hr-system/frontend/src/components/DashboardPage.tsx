import React from "react";
import { SectionPage } from "./SectionPage";

/**
 * DashboardPage serves as the landing page for authenticated users. It
 * displays a friendly overview message and will later be extended with
 * charts and metrics as the HR system matures.
 */
export function DashboardPage() {
    return (
        <SectionPage
            title="Dashboard"
            description="Overview and high‑level HR metrics will be displayed here."
        />
    );
}
