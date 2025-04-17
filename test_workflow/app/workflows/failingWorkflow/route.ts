import { serve } from "@upstash/workflow/nextjs";

export const { POST } = serve(
  async (context) => {

    throw new Error("failed to execute step")
    
    await context.run("step-1", async () => {
      console.log("step-1");
      return "step-1-output";
    });
    
  }, {
    retries: 10,
  }
);
