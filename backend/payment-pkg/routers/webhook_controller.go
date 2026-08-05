package routers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v86/webhook"
)

func (c *controller) IntentUpdate(w http.ResponseWriter, r *http.Request) {
	strWebHScreeet := os.Getenv("PAY_STRIPE_WEBHOOK_SECRET")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	stripeEvent, err := webhook.ConstructEvent(
		body,
		r.Header.Get("Stripe-Signature"),
		strWebHScreeet,
	)
	if err != nil {
		fmt.Println("routers.StripeWebhook:", err.Error())
		w.WriteHeader(http.StatusNotFound)
		//http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if err = c.paymentService.UpdatePaymentStatus(r.Context(), stripeEvent); err != nil {
		fmt.Println("routers.StripeWebhook:", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//switch stripeEvent.Type {
	//case "payment_intent.succeeded":
	//	fmt.Println("webhook event succeeded")
	//	// TODO:
	//	// update payment table
	//	// update order status
	//	// consume reservation
	//case "payment_intent.payment_failed":
	//	fmt.Println("webhook event failed")
	////	// TODO:
	////	// payment = failed
	//case "payment_intent.canceled":
	////	// TODO:
	////	// release reservation
	//default:
	//	// ignore other events
	//}
	w.WriteHeader(http.StatusOK)
	return
}
