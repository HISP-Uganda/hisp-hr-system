import React from "react";
import { Paper, Stack, Typography } from "@mui/material";

type SectionPageProps = {
  title: string;
  description: string;
};

export function SectionPage({ title, description }: SectionPageProps) {
  return (
    <Paper elevation={0} sx={{ p: 4, borderRadius: 2.5, border: "1px solid #dbe3ef", backgroundColor: "rgba(255,255,255,0.88)" }}>
      <Stack spacing={1.2}>
        <Typography variant="h4" sx={{ fontWeight: 800 }}>{title}</Typography>
        <Typography variant="body1" color="text.secondary">{description}</Typography>
      </Stack>
    </Paper>
  );
}
