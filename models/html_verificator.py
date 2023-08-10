# -*- coding: utf-8 -*-
"""HTML Verificator.ipynb

Automatically generated by Colaboratory.

Original file is located at
    https://colab.research.google.com/drive/1ACw4ewe70LE2svjJT_f6nHK2yggo3Ayn

# Phishing website detection based on HTML content
> We are using a basic text classifier (RNN) to determine if the website is a phishing one, or legitm one.

Import required modules
"""

import numpy as np

import tensorflow_datasets as tfds
import tensorflow as tf

tf.get_logger().setLevel('ERROR')

"""Set basic variables"""

AUTOTUNE = tf.data.AUTOTUNE
batch_size = 64
seed = 42

"""Set dataset"""

train_dataset = tf.keras.utils.text_dataset_from_directory(
    'drive/MyDrive/is_ai',
    batch_size=batch_size,
    validation_split=0.2,
    subset='training',
    seed=seed)

class_names = train_dataset.class_names
train_ds = train_dataset.cache().prefetch(buffer_size=AUTOTUNE)

test_dataset = tf.keras.utils.text_dataset_from_directory(
    'drive/MyDrive/is_ai',
    batch_size=batch_size,
    validation_split=0.2,
    subset='validation',
    seed=seed)

test_dataset = test_dataset.cache().prefetch(buffer_size=AUTOTUNE)

for example, label in train_dataset.take(1):
  print('text: ', example.numpy()[:2])
  print('label: ', label.numpy()[:2])

"""Create the text encoder"""

VOCAB_SIZE = 5000
encoder = tf.keras.layers.TextVectorization(
    max_tokens=VOCAB_SIZE)
encoder.adapt(train_dataset.map(lambda text, label: text))

"""Create model"""

model = tf.keras.Sequential([
    encoder,
    tf.keras.layers.Embedding(
        input_dim=len(encoder.get_vocabulary()),
        output_dim=64,
        # Use masking to handle the variable sequence lengths
        mask_zero=True),
    tf.keras.layers.Bidirectional(tf.keras.layers.LSTM(64)),
    tf.keras.layers.Dense(64, activation='relu'),
    tf.keras.layers.Dense(1)
])

"""Compile model to start training"""

model.compile(loss=tf.keras.losses.BinaryCrossentropy(from_logits=True),
              optimizer=tf.keras.optimizers.Adam(1e-4),
              metrics=['accuracy'])

"""Train model"""

history = model.fit(train_dataset, epochs=10,
                    validation_data=test_dataset,
                    validation_steps=30)

test_loss, test_acc = model.evaluate(test_dataset)

print('Test Loss:', test_loss)
print('Test Accuracy:', test_acc)

"""Stack second model"""

model = tf.keras.Sequential([
    encoder,
    tf.keras.layers.Embedding(len(encoder.get_vocabulary()), 64, mask_zero=True),
    tf.keras.layers.Bidirectional(tf.keras.layers.LSTM(64,  return_sequences=True)),
    tf.keras.layers.Bidirectional(tf.keras.layers.LSTM(32)),
    tf.keras.layers.Dense(64, activation='relu'),
    tf.keras.layers.Dropout(0.5),
    tf.keras.layers.Dense(1)
])

model.compile(loss=tf.keras.losses.BinaryCrossentropy(from_logits=True),
              optimizer=tf.keras.optimizers.Adam(1e-4),
              metrics=['accuracy'])

history = model.fit(train_dataset, epochs=10,
                    validation_data=test_dataset,
                    validation_steps=30)

test_loss, test_acc = model.evaluate(test_dataset)

print('Test Loss:', test_loss)
print('Test Accuracy:', test_acc)

"""Use model"""

text_text = ""

for text, _ in train_dataset.take(1):
    text_text = text[:1]

predictions = model.predict(text_text)
print(predictions)

model.save("phishing")

"""Use model"""

new_model = tf.keras.models.load_model("phishing")

predictions = new_model.predict(text_text)
print(predictions)
