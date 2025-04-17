import { serve } from "@upstash/workflow/nextjs";

export const { POST } = serve(
  async (context) => {
    await context.run("step-1", async () => {
      console.log("step-1");
      return "step-1-output";
    });

    await context.run("step-2", async () => {
      console.log("step-2");
      return "step-2-output";
    });

    await context.run("step-3", () => {
      console.log("step-3");
      return "step-3-output";
    });
  }, {
    verbose: true
  }
);
