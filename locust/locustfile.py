from locust import FastHttpUser, task, between
import random

class ProductUser(FastHttpUser):
    wait_time = between(1, 2)

    @task(3)
    def get_product(self):
        self.client.get("/products/1", name="GET /products/[id]")

    @task(1)
    def add_product(self):
        pid = random.randint(1, 1000)
        self.client.post(
            f"/products/{pid}/details",
            name="POST /products/[id]/details",
            json={
                "product_id": pid,
                "sku": "ABC-123-XYZ",
                "manufacturer": "Acme Corporation",
                "category_id": 456,
                "weight": 1250,
                "some_other_id": 789,
            },
        )
