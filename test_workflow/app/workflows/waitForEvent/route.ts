import { serve } from "@upstash/workflow/nextjs";

interface WaitForEventWorkflowRequest {
    eventId: string;
    expectedEventData: any;
}

export const { POST } = serve<WaitForEventWorkflowRequest>(
  async (context) => {
    const requestPayload = context.requestPayload

    console.log(requestPayload)
    await context.run("step-1", async () => {
      console.log("step-1");
      return "step-1-output";
    });

    const {eventData, timeout} = await context.waitForEvent("wait-for-event", requestPayload.eventId, {
        timeout: "24h",
    })
    
    console.log("1", JSON.stringify(requestPayload.expectedEventData))
    console.log("2", JSON.stringify(eventData))
    if (JSON.stringify(requestPayload.expectedEventData) != JSON.stringify(eventData)) {
        throw new Error("failed to execute step")
    }

    await context.run("step-3", () => {
      console.log("step-3");
      return "step-3-output";
    });
  }, {
    retries: 0
  }
);
